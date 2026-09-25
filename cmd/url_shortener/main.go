package main

import (
	"URLShortener/internal/config"
	redirecthandler "URLShortener/internal/http/handlers/redirect"
	urlhandler "URLShortener/internal/http/handlers/url"
	"URLShortener/internal/service"
	"URLShortener/internal/storage/postgres"
	"fmt"
	"log"
	"net/http"

	_ "URLShortener/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title			URLShortener API
// @version		1.0
// @description	A URL shortening service API.
// @host			localhost:8080
// @BasePath		/
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	urlStorage, err := postgres.NewStorage(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := urlStorage.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	urlService := service.NewService(urlStorage)

	urlHandler := urlhandler.NewHandler(urlService)
	redirectHandler := redirecthandler.NewHandler(urlService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /url", urlHandler.CreateURL)
	mux.HandleFunc("GET /urls", urlHandler.AllURLs)
	mux.HandleFunc("GET /{alias}", redirectHandler.Redirect)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Starting server on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
