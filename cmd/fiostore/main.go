package main

import (
	fiostore "github.com/dapixio/fio-shopping"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/send_request", fiostore.ReqHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
