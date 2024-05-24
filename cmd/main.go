package main

import (
	application "beli-mang"
	"beli-mang/config"
	"context"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(err)
	}

	application.Start(cfg)
}
