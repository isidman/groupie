// handlers.go
// Περιέχει τους HTTP handlers του server.
// Κάθε handler αντιστοιχεί σε ένα route.
//
// Οι handlers αντλούν δεδομένα από το api package
// και τα περνάνε στα HTML templates.

package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"groupie-tracker/internal"
)

type Handlers struct {
	api      *internal.Client
	renderer *Renderer
}

func NewHandlers(apiClient *internal.Client, renderer *Renderer) *Handlers {
	return &Handlers{api: apiClient, renderer: renderer}
}

// GET /
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w)
		return
	}
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w)
		return
	}

	artists, err := h.api.GetArtists()
	if err != nil {
		h.serverError(w, err)
		return
	}

	if err := h.renderer.Render(w, "index.html", artists); err != nil {
		h.serverError(w, err)
		return
	}
}

// GET /artist?id=1
func (h *Handlers) Artist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.badRequest(w, "Λείπει το query parameter 'id'.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.badRequest(w, "Το 'id' πρέπει να είναι θετικός ακέραιος.")
		return
	}

	artists, err := h.api.GetArtists()
	if err != nil {
		h.serverError(w, err)
		return
	}

	var found *internal.Artist
	for i := range artists {
		if artists[i].ID == id {
			found = &artists[i]
			break
		}
	}
	if found == nil {
		h.notFound(w)
		return
	}

	rel, err := h.api.GetRelationByID(id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	pageData := internal.BuildArtistPageData(*found, rel)

	if err := h.renderer.Render(w, "artist.html", pageData); err != nil {
		h.serverError(w, err)
		return
	}
}

// GET /api/artist?id=1  (JSON response - για event/action αργότερα)
func (h *Handlers) ArtistAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "invalid id"})
		return
	}

	artists, err := h.api.GetArtists()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "api error"})
		return
	}

	var found *internal.Artist
	for i := range artists {
		if artists[i].ID == id {
			found = &artists[i]
			break
		}
	}
	if found == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "not found"})
		return
	}

	rel, err := h.api.GetRelationByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "relation error"})
		return
	}

	data := internal.BuildArtistPageData(*found, rel)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(data)
}
