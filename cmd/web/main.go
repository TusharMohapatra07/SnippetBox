package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"snippetbox/internal/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error(err.Error())
	}
	connStr := os.Getenv("CONNSTR")
	addr := os.Getenv("ADDRESS")

	db, err := openDB(connStr)
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	mux := app.routes()

	logger.Info("Starting server", "addr", addr)
	err = http.ListenAndServe(addr, mux)

	logger.Error(err.Error())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
