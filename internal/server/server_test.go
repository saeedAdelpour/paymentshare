package server_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/saeedAdelpour/paymentshare/internal/server"
)

func TestStartupRejectsInvalidDatabaseConfigurationWithoutExposingSecrets(t *testing.T) {
	var logs bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := server.Run(ctx, server.Config{
		HTTPAddr:    "127.0.0.1:0",
		DatabaseURL: "postgres://user:secret-marker@host:invalid/database",
	}, slog.New(slog.NewJSONHandler(&logs, nil)))
	if err == nil {
		t.Fatal("startup accepted an invalid database configuration")
	}
	if errors.Unwrap(err) == nil {
		t.Fatal("startup discarded the underlying diagnostic cause")
	}
	if strings.Contains(err.Error()+logs.String(), "secret-marker") {
		t.Fatal("startup exposed database credentials")
	}
}
