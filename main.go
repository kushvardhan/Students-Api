package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){
	r := mux.NewRouter();

	log.Println("Starting Server at PORT:4444.");
	log.Fatal(http.ListenAndServe(":4444",r))
}

