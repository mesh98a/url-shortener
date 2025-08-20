package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/sql"

	mwLogger "url-shortener/internal/http-server/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	fmt.Println("%#v\n", cfg)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug info")

	user := cfg.MYSQL.User
	password := cfg.MYSQL.Password
	host := cfg.MYSQL.Host
	port := cfg.MYSQL.Port

	path := user + ":" + password + "@tcp(" + host + ":" + port + ")/"
	storage, err := sql.NewStorage(path)
	if err != nil {
		slog.Error("storage failed ", err)
		os.Exit(1)
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
