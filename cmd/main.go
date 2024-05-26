package main

import (
	application "beli-mang"
	"beli-mang/config"
	"context"
	// see https://pkg.go.dev/net/http/pprof for the docs
	_ "net/http/pprof"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

	application.Start(cfg)
}
