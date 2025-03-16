package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"
	"toramanomer/snippetbox/internal/models"

	"net/http"

	"github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func openDB(driverName string, dataSourceName string) (*sql.DB, error) {

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {

	var (
		logger   = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		mySqlCfg = &mysql.Config{
			User:      "web",
			Passwd:    "pass",
			DBName:    "snippetbox",
			Net:       "tcp",
			ParseTime: true,
		}

		driverName     = "mysql"
		dataSourceName = mySqlCfg.FormatDSN()
		db, err        = openDB(driverName, dataSourceName)

		cfg = config{}
	)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Close connection program exits
	defer db.Close()

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")

	flag.Parse()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	server := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", slog.String("addr", server.Addr))
	app.logger.Error(server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem").Error())
}
