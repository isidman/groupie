// merge.go
// Περιέχει τη λογική που ενώνει δεδομένα
// από διαφορετικά endpoints του API
// σε μία δομή που είναι εύχρηστη για τα HTML templates.
//
// Δεν πραγματοποιεί HTTP requests.

package internal

import (
	"sort"
	"strings"
)

// BuildArtistPageData builds all page data for the artist page
func BuildArtistPageData(a Artist, r Relation, l Locations, d Dates) ArtistPageData {
	concerts := make([]Concert, 0, len(r.DatesLocations))

	for rawLoc, dates := range r.DatesLocations {
		location := FormatLocation(rawLoc)

		cleanDates := make([]string, 0, len(dates))
		for _, date := range dates {
			cleanDates = append(cleanDates, strings.TrimPrefix(date, "*"))
		}

		sort.Strings(cleanDates)

		concerts = append(concerts, Concert{
			Location: location,
			Dates:    cleanDates,
		})
	}

	sort.Slice(concerts, func(i, j int) bool {
		return concerts[i].Location < concerts[j].Location
	})

	allLocations := make([]string, 0, len(l.Locations))
	for _, loc := range l.Locations {
		allLocations = append(allLocations, FormatLocation(loc))
	}
	sort.Strings(allLocations)

	allDates := make([]string, 0, len(d.Dates))
	for _, date := range d.Dates {
		allDates = append(allDates, strings.TrimPrefix(date, "*"))
	}
	sort.Strings(allDates)

	rawLocations := make([]string, 0, len(l.Locations))
	for _, loc := range l.Locations {
		rawLocations = append(rawLocations, loc)
	}

	return ArtistPageData{
		Artist:       a,
		Concerts:     concerts,
		AllLocations: allLocations,
		AllDates:     allDates,
		RawLocations: rawLocations,
	}
}

// FormatLocation converts "new_york-usa" to "New york, Usa"
func FormatLocation(raw string) string {
	raw = strings.ReplaceAll(raw, "_", " ")
	parts := strings.Split(raw, "-")

	for i := range parts {
		parts[i] = capitalize(parts[i])
	}

	if len(parts) >= 2 {
		return parts[0] + ", " + strings.Join(parts[1:], ", ")
	}

	return capitalize(raw)
}

func capitalize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
