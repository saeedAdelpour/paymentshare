//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

// This TCP proxy still executes every query against real PostgreSQL. It withholds
// only the unlock response, reproducing a stalled network during Goose cleanup.
func stalledUnlockProxy(t *testing.T, dsn string) (string, <-chan struct{}) {
	t.Helper()
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal("parse test database configuration failed")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	seen := make(chan struct{}, 1)
	var peers sync.Map
	var workers sync.WaitGroup
	accepted := make(chan struct{})
	go func() {
		defer close(accepted)
		for {
			client, err := listener.Accept()
			if err != nil {
				return
			}
			peers.Store(client, true)
			workers.Add(1)
			go func() {
				defer workers.Done()
				defer peers.Delete(client)
				defer func() { _ = client.Close() }()
				upstream, err := net.DialTimeout("tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port))), 2*time.Second)
				if err != nil {
					return
				}
				defer func() { _ = upstream.Close() }()
				monitor := &unlockMonitor{destination: upstream, seen: seen}
				copied := make(chan struct{})
				go func() {
					defer close(copied)
					_, _ = io.Copy(monitor, client)
					_ = upstream.Close()
				}()
				_, _ = io.Copy(&unlockResponse{destination: client, monitor: monitor}, upstream)
				_ = client.Close()
				<-copied
			}()
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		<-accepted
		peers.Range(func(key, _ any) bool { _ = key.(net.Conn).Close(); return true })
		workers.Wait()
	})
	// The cancellation fault injection runs on the dedicated local TCP test DB.
	proxy := &url.URL{Scheme: "postgres", Host: listener.Addr().String(), Path: "/" + cfg.Database, User: url.UserPassword(cfg.User, cfg.Password)}
	query := url.Values{"sslmode": {"disable"}}
	for name, value := range cfg.RuntimeParams {
		query.Set(name, value)
	}
	proxy.RawQuery = query.Encode()
	return proxy.String(), seen
}

type unlockMonitor struct {
	destination io.Writer
	seen        chan<- struct{}
	stalled     atomic.Bool
	tail        []byte
}

func (m *unlockMonitor) Write(data []byte) (int, error) {
	const marker = "pg_advisory_unlock"
	combined := append(m.tail, data...)
	if bytes.Contains(combined, []byte(marker)) && m.stalled.CompareAndSwap(false, true) {
		select {
		case m.seen <- struct{}{}:
		default:
		}
	}
	if len(combined) > len(marker) {
		combined = combined[len(combined)-len(marker):]
	}
	m.tail = append([]byte(nil), combined...)
	return m.destination.Write(data)
}

type unlockResponse struct {
	destination io.Writer
	monitor     *unlockMonitor
}

func (w *unlockResponse) Write(data []byte) (int, error) {
	if w.monitor.stalled.Load() {
		return len(data), nil
	}
	return w.destination.Write(data)
}

func TestMigrationCancellationBoundsUnlockCleanup(t *testing.T) {
	dsn, _ := newDatabase(t)
	proxy, unlockStarted := stalledUnlockProxy(t, dsn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, migrationBinary, "-to", "1")
	command.Env = append(os.Environ(), "DATABASE_URL="+proxy)
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- command.Wait() }()
	select {
	case <-unlockStarted:
	case err := <-done:
		t.Fatalf("migration exited before unlock fault: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("migration did not reach unlock")
	}
	if err := command.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("migration succeeded despite a failed unlock")
		}
		if ctx.Err() != nil {
			t.Fatal("migration needed forced termination")
		}
	case <-time.After(7 * time.Second):
		_ = command.Process.Kill()
		if err := <-done; err == nil {
			t.Error("forced termination unexpectedly succeeded")
		}
		t.Fatal("migration unlock exceeded the bounded cleanup allowance")
	}
	runMigration(t, dsn, "1", true)
}
