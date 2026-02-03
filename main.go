// main.go
// main αρχικοποιεί τον server και ξεκινά να ακούει εισερχόμενα αιτήματα.
// Αν αποτύχει η εκκίνηση του server, η εφαρμογή τερματίζεται με σφάλμα.

// Σημείο εκκίνησης της εφαρμογής.
// Εδώ γίνεται η αρχικοποίηση του HTTP server,
// ο ορισμός των routes και η εκκίνηση του ListenAndServe.
//
// Δεν περιέχει επιχειρησιακή λογική (business logic),
// μόνο συντονισμό (routes, static αρχεία, handlers).

package main

import (
	"log"
	"net/http"
	"time"

	"groupie-tracker/internal"
	"groupie-tracker/web"
)

func main() {
	// Client για το εξωτερικό API με timeout (ώστε να μην "κολλάει" ποτέ)
	apiClient := internal.NewClient(10 * time.Second)

	// Renderer για τα templates (εσύ έχεις τα html μέσα στον φάκελο web/)
	renderer, err := web.NewRenderer("web/*.html")
	if err != nil {
		log.Fatal("Αποτυχία φόρτωσης templates: ", err)
	}

	// Handlers
	h := web.NewHandlers(apiClient, renderer)

	mux := http.NewServeMux()

	// Routes σελίδων
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/artist", h.Artist)

	// Route JSON (για event/action αργότερα)
	mux.HandleFunc("/api/artist", h.ArtistAPI)

	// Static αρχεία (CSS/JS) -> εσύ έχεις φάκελο "static/"
	staticFS := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	// Middleware: logging + recover (να μην κρασάρει ποτέ)
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
