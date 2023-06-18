package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// holds application-wide dependencies
type application struct {
	infoLog *log.Logger
	errLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "provide addr to run server")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application { infoLog, errLog }

	mux := app.routes()

	srv := http.Server {
		Addr: *addr,
		Handler: mux,
		ErrorLog: errLog,
	}

	infoLog.Println("Starting at ", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errLog.Fatal(err)
	}
}
