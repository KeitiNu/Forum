package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"git.01.kood.tech/roosarula/forum/pkg/data"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	models        data.Models
	templateCache map[string]*template.Template
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	err = Migrate(db)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Data and configuration for app
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		models:        data.NewModels(db),
		templateCache: templateCache,
	}
	// We create a http.Server configuration
	srv := &http.Server{
		Addr:     ":8090",
		ErrorLog: errorLog,
		Handler:  app.routes(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("Started server on http://localhost%s\n", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

// The openDB() function returns a sql.DB connection pool.
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:database.db")
	if err != nil {
		return nil, err
	}

	// Use Ping to establish a new connection to the database, if the connection couldn't be
	// established successfully this will return an
	// error.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
