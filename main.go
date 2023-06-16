package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from gomibako"))
}

// handler to showing individual gomi
func showGomi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("handler for showing gomi"))
}

// handler to create new gomi
func createGomi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not supported", 405)
		return
	}

	w.Write([]byte("create gomi"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/gomi", showGomi)
	mux.HandleFunc("/gomi/create", createGomi)

	log.Println("Starting at :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
