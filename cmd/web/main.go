package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/MichaelSWong/snippetbox/internal/models"
	"github.com/MichaelSWong/snippetbox/internal/store"
	"github.com/MichaelSWong/snippetbox/migrations"
	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable", "Postgres data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := store.OpenDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = store.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
