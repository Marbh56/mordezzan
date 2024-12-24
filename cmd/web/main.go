package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Info("error loading .env file")
		os.Exit(1)
	}

	connStr := os.Getenv("dbstring")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Info("error opening db")
		os.Exit(1)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Info("cant ping db")
		os.Exit(1)
	}

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
