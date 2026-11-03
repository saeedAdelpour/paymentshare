package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/saeedAdelpour/paymentshare/internal/contract"
)

func main() {
	path := flag.String("file", "api/openapi.yaml", "OpenAPI document to validate")
	flag.Parse()
	document, err := os.ReadFile(*path)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = contract.Load(ctx, document)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "OpenAPI validation failed:", err)
		os.Exit(1)
	}
	fmt.Println("OpenAPI contract is valid")
}
