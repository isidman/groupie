// models.go
// Περιέχει όλα τα structs που αντιστοιχούν
// στα JSON responses του εξωτερικού Groupie Trackers API.
//
// Τα structs αυτά χρησιμοποιούνται μόνο για αποθήκευση δεδομένων
// και όχι για απευθείας εμφάνιση στα templates.

package internal

// Artist represents one artist/band from /api/artists
type Artist struct {
	ID               int      `json:"id"`
	Image            string   `json:"image"`
	Name             string   `json:"name"`
	Members          []string `json:"members"`
	CreationDate     int      `json:"creationDate"`
	FirstAlbum       string   `json:"firstAlbum"`
	LocationsURL     string   `json:"locations"`
	ConcertLocations []string `json:"-"`
	DatesURL         string   `json:"concertDates"`
	RelationsURL     string   `json:"relations"`
}

// Relation represents /api/relation/:id
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Locations represents /api/locations/:id
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

// Dates represents /api/dates/:id
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Concert struct {
	Location string
	Dates    []string
}

type ArtistPageData struct {
	Artist       Artist
	Concerts     []Concert
	AllLocations []string
	AllDates     []string
	RawLocations []string
}

type SearchResult struct {
	ArtistID   int    `json:"artistID"`
	ArtistName string `json:"artistName"`
	Match      string `json:"match"`
	Type       string `json:"type"`
}

type GeoLocation struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}
