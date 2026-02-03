// merge.go
// Περιέχει τη λογική που ενώνει δεδομένα
// από διαφορετικά endpoints του API
// σε μία δομή που είναι εύχρηστη για τα HTML templates.
//
// Δεν πραγματοποιεί HTTP κλήσεις.

package internal

import (
	"sort"
	"strings"
)

// BuildArtistPageData ενώνει Artist + Relation σε κάτι έτοιμο για template
func BuildArtistPageData(a Artist, r Relation) ArtistPageData {
	concerts := make([]Concert, 0, len(r.DatesLocations))

	for rawLoc, dates := range r.DatesLocations {
		loc := FormatLocation(rawLoc)

		cleanDates := make([]string, 0, len(dates))
		for _, d := range dates {
			cleanDates = append(cleanDates, strings.TrimPrefix(d, "*"))
		}

		sort.Strings(cleanDates)
		concerts = append(concerts, Concert{
			Location: loc,
			Dates:    cleanDates,
		})
	}

	// Σταθερή σειρά (alphabetical) για ωραία εμφάνιση
	sort.Slice(concerts, func(i, j int) bool {
		return concerts[i].Location < concerts[j].Location
	})

	return ArtistPageData{
		Artist:   a,
		Concerts: concerts,
	}
}

// FormatLocation μετατρέπει "athens-greece" ή "new_york-usa" σε "Athens, Greece"
func FormatLocation(raw string) string {
	raw = strings.ReplaceAll(raw, "_", " ")
	parts := strings.Split(raw, "-")
	for i := range parts {
		parts[i] = title(parts[i])
	}
	if len(parts) >= 2 {
		return parts[0] + ", " + strings.Join(parts[1:], ", ")
	}
	return title(raw)
}

func title(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
