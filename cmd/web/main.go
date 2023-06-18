package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

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

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/gomi", app.showGomi)
	mux.HandleFunc("/gomi/create", app.createGomi)

	// file server to serve static files
	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

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
