package main

import (
	redirecthandler "URLShortener/internal/http/handlers/redirect"
	urlhandler "URLShortener/internal/http/handlers/url"
	"URLShortener/internal/service"
	"URLShortener/internal/storage/memory"
	"fmt"
	"net/http"
)

func main() {
	urlStorage := memory.NewStorage()
	urlService := service.NewService(urlStorage)
	urlHandler := urlhandler.NewHandler(urlService)
	redirectHandler := redirecthandler.NewHandler(urlService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /url", urlHandler.CreateURL)
	mux.HandleFunc("GET /urls", urlHandler.AllURLs)
	mux.HandleFunc("GET /{alias}", redirectHandler.Redirect)

	fmt.Println("Starting server on http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
