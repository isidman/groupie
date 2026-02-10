package web

import (
	"log"
	"net/http"
)

func (h *Handlers) badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Bad Request",
		"Message": msg,
	})
}

func (h *Handlers) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Not Found",
		"Message": "The page you requested does not exist.",
	})
}

func (h *Handlers) methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Method Not Allowed",
		"Message": "This HTTP method is not allowed here.",
	})
}

func (h *Handlers) serverError(w http.ResponseWriter, err error) {
	log.Println("server error:", err)
	w.WriteHeader(http.StatusInternalServerError)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Server Error",
		"Message": "Something went wrong. Please try again.",
	})
}
