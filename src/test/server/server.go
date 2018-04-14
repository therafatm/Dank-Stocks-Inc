package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func respond(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	w.Write("ok")
}

func main() {
	router := mux.NewRouter()
	port := os.Getenv("TRANS_PORT")
	router.HandleFunc("/*", respond)
	server := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  2 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	log.Println("Running transaction server on port: " + port)
	log.Fatal(server.ListenAndServe())
}