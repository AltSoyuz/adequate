package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AltSoyuz/adequate/internal/migration"
	"github.com/AltSoyuz/adequate/internal/store"
	"github.com/AltSoyuz/adequate/lib/buildinfo"
	"github.com/AltSoyuz/adequate/lib/envflag"
	"github.com/AltSoyuz/adequate/lib/httpserver"
	"github.com/AltSoyuz/adequate/lib/logger"
)

var (
	httpAddr      = flag.String("http.listenAddr", ":8080", "HTTP listen address")
	sqlitePath    = flag.String("store.sqlitePath", "data/db", "SQLite database file path")
	staticDirPath = flag.String("http.staticDir", "", "Static files directory (for serving UI assets)")
)

func main() {
	envflag.Parse()
	logger.Init()
	buildinfo.Init()
	startime := time.Now()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	store := store.Init(ctx, *sqlitePath)
	defer store.Close()

	mux := http.NewServeMux()

	addRoutes(mux, store)

	if *staticDirPath != "" {
		logger.Info("ui app", "prefix", "/", "staticDir", *staticDirPath)
		mux.HandleFunc("/", httpserver.SPAFileServer(*staticDirPath))
	}

	logger.Info("started app", "duration", time.Since(startime).String())

	if err := httpserver.Serve(ctx, *httpAddr, mux); err != nil {
		logger.Fatal("http serve", "err", err)
	}

	logger.Info("graceful shutdown completed")
}

func addRoutes(mux *http.ServeMux, store *store.Store) {
	mux.HandleFunc("/api/migrations/version", migration.MigrationHandler(store))
}
