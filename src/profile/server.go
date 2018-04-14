package main

import (
	"log"
	"os"
	"time"
	"net/http"
	"github.com/gorilla/mux"
)

func respond(w http.ResponseWriter, r *http.Request){
	log.Println(r)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func main() {
	router := mux.NewRouter()
	port := os.Getenv("TRANS_PORT")
	router.HandleFunc("/", respond)
	router.HandleFunc("/api/clearUsers",respond)
	router.HandleFunc("/api/availableBalance/{username}/{trans}", respond)
	router.HandleFunc("/api/availableShares/{username}/{symbol}/{trans}", respond)
	router.HandleFunc("/api/add/{username}/{money}/{trans}", respond)
	router.HandleFunc("/api/getQuote/{username}/{symbol}/{trans}", respond)
	router.HandleFunc("/api/buy/{username}/{symbol}/{amount}/{trans}", respond)
	router.HandleFunc("/api/commitBuy/{username}/{trans}", respond)
	router.HandleFunc("/api/cancelBuy/{username}/{trans}", respond)
	router.HandleFunc("/api/sell/{username}/{symbol}/{amount}/{trans}", respond)
	router.HandleFunc("/api/commitSell/{username}/{trans}", respond)
	router.HandleFunc("/api/cancelSell/{username}/{trans}", respond)
	router.HandleFunc("/api/setBuyAmount/{username}/{symbol}/{amount}/{trans}", respond)
	router.HandleFunc("/api/setBuyTrigger/{username}/{symbol}/{triggerPrice}/{trans}", respond)
	router.HandleFunc("/api/cancelSetBuy/{username}/{symbol}/{trans}", respond)
	router.HandleFunc("/api/setSellAmount/{username}/{symbol}/{amount}/{trans}", respond)
	router.HandleFunc("/api/cancelSetSell/{username}/{symbol}/{trans}", respond)
	router.HandleFunc("/api/setSellTrigger/{username}/{symbol}/{triggerPrice}/{trans}", respond)
	router.HandleFunc("/api/dumplog/{filename}/{trans}", respond)
	router.HandleFunc("/api/dumplog/{filename}/{username}/{trans}", respond)
	router.HandleFunc("/api/displaySummary/{username}/{trans}", respond)
	server := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  2 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	log.Println("Running profile server on port: " + port)	
	err := server.ListenAndServe()
	if err != nil {
        log.Fatal("ListenAndServe: ", err)
	}
}