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

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
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

	//DATABASE
	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	err = RunMigrateScripts(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully migrated DB..")

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

//Run migrate scripts to create database if not created before.
func RunMigrateScripts(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite3 db driver failed %s", err)
	}

	res := bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		})

	d, err := bindata.WithInstance(res)
	m, err := migrate.NewWithInstance("go-bindata", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrating database failed %s", err)
	}

	return nil
}
