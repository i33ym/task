package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"web.taswiya-todo.cc/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *log.Logger
	models *models.Models
}


func main() {
	logger := log.New(os.Stdout, "[ToDo]\t", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Println("Kicking off...")

	var addr string
	var dsn string
	flag.StringVar(&addr, "addr", ":5050", "HTTP address")
	flag.StringVar(&dsn, "dsn", os.Getenv("TASWIYA_TODO_DSN"), "MySQL DSN")
	flag.Parse()

	db, err := open("mysql", dsn)
	if err != nil {
		logger.Printf("failed to open db connection: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Println("Connected to MySQL")
	
	models := models.NewModels(db)
	app := &application{
		logger: logger,
		models: models,
	}

	server := &http.Server{
		Addr: addr,
		Handler: app.routes(),
		ErrorLog: logger,
	}

	logger.Printf("Listening on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Printf("failed to listen on %s: %s", addr, err)
		os.Exit(1)
	}
}

func open(driver string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}