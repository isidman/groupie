package main

import (
	"log"
	"net/http"
	"time"

	"groupie-tracker/internal"
	"groupie-tracker/web"
)

func main() {
	apiClient := internal.NewClient(10 * time.Second)

	cache := internal.NewArtistCache(apiClient)
	cache.Start()

	renderer, err := web.NewRenderer("web/*.html")
	if err != nil {
		log.Fatal("Template loading failed:", err)
	}

	h := web.NewHandlers(apiClient, cache, renderer)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/artist", h.Artist)
	mux.HandleFunc("/api/artist", h.ArtistAPI)
	mux.HandleFunc("/api/search", h.Search)
	mux.HandleFunc("/api/filter", h.Filter)
	mux.HandleFunc("/api/geocode", h.Geocode)

	staticFS := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	handler := web.RecoverMiddleware(web.LoggingMiddleware(mux))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 12 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server running: http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
