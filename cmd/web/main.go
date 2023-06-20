package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abhaysp95/gomibako/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// holds application-wide dependencies
type application struct {
	infoLog *log.Logger
	errLog *log.Logger
	gomi *mysql.GomiModel
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application {
		infoLog: infoLog,
		errLog: errLog,
		gomi: &mysql.GomiModel{  // dependency injection
			DB: nil,
		},
	}
	err := godotenv.Load()
	if err != nil {
		app.errLog.Println("Error loading .env file")
	}

	var addr string
	var dsn string

	flag.StringVar(&addr, "addr", "", "Provide addr to run server")
	flag.StringVar(&dsn, "dsn", "", "Mariadb data source name")
	flag.Parse()

	if (addr == "") {
		addr = os.Getenv("ADDR")
	}
	if (dsn == "") {
		dsn = fmt.Sprintf("%s:%s@/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DB"))
	}

	if (addr == "" || dsn == "") {
		app.errLog.Fatal("Provide initial server and db configuration to run application. See -h for more info")
	}

	db, err := app.openDB(dsn)
	if err != nil {
		app.errLog.Fatal(err)
	}

	defer db.Close()
	app.gomi.DB = db  // assign db pool

	mux := app.routes()

	srv := http.Server {
		Addr: addr,
		Handler: mux,
		ErrorLog: errLog,
	}

	infoLog.Println("Starting at ", addr)
	err = srv.ListenAndServe()
	if err != nil {
		errLog.Fatal(err)
	}
}

func (app *application) openDB(dsn string) (*sql.DB, error)  {
	db, err := sql.Open("mysql", dsn + "?parseTime=true")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	app.infoLog.Println("DB connection successful")

	return db, nil
}
