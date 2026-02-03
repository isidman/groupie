package web

import (
	"log"
	"net/http"
)

func (h *Handlers) badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Λάθος αίτημα",
		"Message": msg,
	})
}

func (h *Handlers) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Δεν βρέθηκε",
		"Message": "Η σελίδα που ζητήσατε δεν υπάρχει.",
	})
}

func (h *Handlers) methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Μη επιτρεπτή μέθοδος",
		"Message": "Η μέθοδος HTTP δεν επιτρέπεται εδώ.",
	})
}

func (h *Handlers) serverError(w http.ResponseWriter, err error) {
	log.Println("server error:", err)
	w.WriteHeader(http.StatusInternalServerError)
	_ = h.renderer.Render(w, "error.html", map[string]any{
		"Title":   "Σφάλμα server",
		"Message": "Κάτι πήγε στραβά. Δοκιμάστε ξανά.",
	})
}
