// handlers.go
// Περιέχει τους HTTP handlers του server.
// Κάθε handler αντιστοιχεί σε ένα route.
//
// Οι handlers αντλούν δεδομένα από το internal package
// και τα περνάνε στα HTML templates.

package web

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/internal"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Handlers συγκεντρώνει τα dependencies που χρειάζονται οι handlers.
type Handlers struct {
	api      *internal.Client
	cache    *internal.ArtistCache
	renderer *Renderer
}

// NewHandlers δημιουργεί ένα νέο Handlers instance.
func NewHandlers(apiClient *internal.Client, cache *internal.ArtistCache, renderer *Renderer) *Handlers {
	return &Handlers{
		api:      apiClient,
		cache:    cache,
		renderer: renderer,
	}
}

// Index handles GET /
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		h.methodNotAllowed(w)
		return
	}

	artists := h.cache.GetArtists()
	if len(artists) == 0 {
		var err error
		artists, err = h.api.GetArtists()
		if err != nil {
			h.serverError(w, err)
			return
		}
	}

	if err := h.renderer.Render(w, "index.html", artists); err != nil {
		h.serverError(w, err)
		return
	}
}

// Artist handles GET /artist?id=1
func (h *Handlers) Artist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.badRequest(w, "Missing query parameter 'id'.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.badRequest(w, "The 'id' must be a positive integer.")
		return
	}

	artists := h.cache.GetArtists()
	if len(artists) == 0 {
		artists, err = h.api.GetArtists()
		if err != nil {
			h.serverError(w, err)
			return
		}
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

	relation, err := h.api.GetRelationByID(id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	locations, err := h.api.GetLocationsByID(id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	dates, err := h.api.GetDatesByID(id)
	if err != nil {
		h.serverError(w, err)
		return
	}

	pageData := internal.BuildArtistPageData(*found, relation, locations, dates)

	if err := h.renderer.Render(w, "artist.html", pageData); err != nil {
		h.serverError(w, err)
		return
	}
}

// ArtistAPI handles GET /api/artist?id=1
func (h *Handlers) ArtistAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Method not allowed.",
		})
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Invalid id.",
		})
		return
	}

	artists := h.cache.GetArtists()
	if len(artists) == 0 {
		artists, err = h.api.GetArtists()
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": "Artists fetch failed.",
			})
			return
		}
	}

	var found *internal.Artist
	for i := range artists {
		if artists[i].ID == id {
			found = &artists[i]
			break
		}
	}

	if found == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Artist not found.",
		})
		return
	}

	relation, err := h.api.GetRelationByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Relation fetch failed.",
		})
		return
	}

	locations, err := h.api.GetLocationsByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Locations fetch failed.",
		})
		return
	}

	dates, err := h.api.GetDatesByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Dates fetch failed.",
		})
		return
	}

	data := internal.BuildArtistPageData(*found, relation, locations, dates)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Method not allowed.",
		})
		return
	}

	query := r.URL.Query().Get("q")

	if len(query) <= 1 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode([]internal.SearchResult{})
		return
	}

	artists := h.cache.GetArtists()
	results := internal.SearchArtists(query, artists)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(results)

}

func (h *Handlers) Filter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Method not allowed.",
		})
		return
	}

	q := r.URL.Query()

	// Read min & max in creation as ints
	creationMin, _ := strconv.Atoi(q.Get("creationMin"))
	creationMax, _ := strconv.Atoi(q.Get("creationMax"))
	albumMin, _ := strconv.Atoi(q.Get("albumMin"))
	albumMax, _ := strconv.Atoi(q.Get("albumMax"))

	// Read members as []int
	// The q["members"] gives back []string
	// It loops through them and convert each to int
	members := []int{}
	for _, m := range q["members"] {
		val, err := strconv.Atoi(m)
		if err == nil {
			members = append(members, val)
		}
	}

	location := q.Get("location")

	//Filter Parameters
	params := internal.FilterParams{
		CreationMin: creationMin,
		CreationMax: creationMax,
		AlbumMin:    albumMin,
		AlbumMax:    albumMax,
		Members:     members,
		Location:    location,
	}

	artists := h.cache.GetArtists()
	results := internal.FilterArtists(artists, params)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(results)
}

func (h *Handlers) Geocode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Method not allowed.",
		})
		return
	}

	location := r.URL.Query().Get("location")

	if location == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Missing location.",
		})
		return
	}

	cleaned := strings.ReplaceAll(location, "_", " ")
	parts := strings.Split(cleaned, "-")

	var query string
	if len(parts) >= 2 && !strings.Contains(parts[len(parts)-1], " ") {
		// Single word country — put it first: "ex:japan osaka"
		country := parts[len(parts)-1]
		city := strings.Join(parts[:len(parts)-1], " ")
		query = country + " " + city
	} else {
		// Multi-word country — keep forward order: "ex: abu dhabi united arab emirates"
		query = strings.Join(parts, " ")
	}

	apiURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1&featuretype=city",
		url.QueryEscape(query),
	)

	// Create and send request - if it doesn't work respond with errors
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Internal server error.",
		})
		return
	}
	req.Header.Set("User-Agent", "groupie-tracker/1.0")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusRequestTimeout)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Request timeout.",
		})
		return
	}
	defer resp.Body.Close()

	//Decoding response
	var results []internal.GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Internal server error.",
		})
		return
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "No results found.",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(results[0])
}
