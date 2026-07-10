package main

import (
	"log"
	"net/http"

	"backend/api"
)

func main() {
	router := api.RegisterRoutes()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
