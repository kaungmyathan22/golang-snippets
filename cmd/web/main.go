package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kaungmyathan22/golang-sinppets/pkg/models/mysql"
)

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")

	dsn := flag.String("dsn", "root:rootpassword@tcp(localhost:3306)/snippetbox?parseTime=true&tls=skip-verify", "MySQL data source name")

	flag.Parse()
	// f, err := os.OpenFile("./logs/info.log", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	errCh := make(chan error, 1)

	// Start the server in a goroutine and check for errors
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			errCh <- err
		}
	}()

	// Wait for the error channel to receive an error or close
	err = <-errCh
	errorLog.Fatal(err) // Log and exit gracefully
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(5)
	infoLog.Println("Successfully connected to the database.")
	return db, nil
}
