package main

import (
	"hsbc/handler"
	"hsbc/services"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	service := services.NewGeoNamesApi(client)
	geoNamesHandler := handler.NewHandler(service)
	
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.Handle("/api/cities/count", geoNamesHandler)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
