// models.go
// Περιέχει όλα τα structs που αντιστοιχούν
// στα JSON responses του εξωτερικού Groupie Trackers API.
//
// Τα structs αυτά χρησιμοποιούνται μόνο για αποθήκευση δεδομένων
// και όχι για απευθείας εμφάνιση στα templates.

package internal

// Artist αναπαριστά έναν καλλιτέχνη/συγκρότημα όπως έρχεται από /api/artists
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocationsURL string   `json:"locations"`
	DatesURL     string   `json:"concertDates"`
	RelationsURL string   `json:"relations"`
}

// Relation αναπαριστά τα datesLocations από /api/relation/:id
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// View model για τη σελίδα artist (εύκολο για templates)
type Concert struct {
	Location string
	Dates    []string
}

type ArtistPageData struct {
	Artist   Artist
	Concerts []Concert
}
