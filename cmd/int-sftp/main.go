package main

import (
	"context"
	"os"

	"devops1/internal/config"
	"devops1/internal/db"
	sftppkg "devops1/internal/sftp"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("component", "sftp").Logger()
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		panic(err)
	}
	pool, err := db.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	server := &sftppkg.Server{Cfg: cfg, Queries: db.New(pool), Logger: logger, UserType: "internal"}
	if err := server.Start(context.Background(), cfg.SFTP.InternalAddress); err != nil {
		panic(err)
	}
}
