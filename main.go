package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	v1 "github.com/billyboar/battleships/api/v1"
)

var portFlag string

func main() {
	// var router = mux.NewRouter()
	// router.HandleFunc("/health", healthCheck).Methods("GET")

	var port int

	flag.IntVar(&port, "port", 3000, "Port number to run server on")

	server, err := v1.NewAPIServer()
	if err != nil {
		panic(err)
	}
	server.RegisterRoutes()

	fmt.Println(fmt.Sprintf("Running server on :%d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server.Router))
}
