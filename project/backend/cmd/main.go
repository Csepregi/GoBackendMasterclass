package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"eats/backend"
	"eats/backend/common/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	log.Init(slog.LevelInfo)

	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		panic("POSTGRES_URL environment variable is not set")
	}

	dbPgx, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	app, err := backend.New(
		ctx,
		dbPgx,
	)
	if err != nil {
		panic(err)
	}

	app.Run(ctx, ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Service initialized")
}
