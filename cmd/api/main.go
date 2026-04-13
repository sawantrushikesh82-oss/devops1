package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"devops1/internal/api"
	"devops1/internal/config"
	"devops1/internal/db"
	"devops1/internal/metrics"
	syncpkg "devops1/internal/sync"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("component", "api").Logger()
	ctx := logger.WithContext(context.Background())
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		panic(err)
	}
	pool, err := db.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	metrics.Register()
	dbsql := stdlib.OpenDBFromPool(pool)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(dbsql, "internal/db/migrations"); err != nil {
		panic(err)
	}
	q := db.New(pool)
	eng := syncpkg.NewEngine(cfg, q, logger)
	go eng.Start(ctx)
	pub, priv := mustKeys(cfg.API.JWTPublicKeyPath, cfg.API.JWTPrivateKeyPath)
	r := api.NewServer(cfg, q, eng, logger, pub, priv)
	addr := fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port)
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

func mustKeys(pubPath, privPath string) (*rsa.PublicKey, *rsa.PrivateKey) {
	pubPem, _ := os.ReadFile(pubPath)
	privPem, _ := os.ReadFile(privPath)
	pb, _ := pem.Decode(pubPem)
	pr, _ := pem.Decode(privPem)
	pubAny, _ := x509.ParsePKIXPublicKey(pb.Bytes)
	priv, _ := x509.ParsePKCS8PrivateKey(pr.Bytes)
	return pubAny.(*rsa.PublicKey), priv.(*rsa.PrivateKey)
}
