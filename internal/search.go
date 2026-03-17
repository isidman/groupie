package internal

import (
	"fmt"
	"strings"
)

// This funky function does the search baby based on your "query"
// and matching it to different types of data (slice of artists).
// It gives back a slice of "SearchResult" babyyyyyyyyyyyyyyyyyyy.
func SearchArtists(query string, artists []Artist) []SearchResult {
	results := []SearchResult{}
	query = strings.ToLower(query)

	for _, artist := range artists {

		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, SearchResult{
				ArtistID:   artist.ID,
				ArtistName: artist.Name,
				Match:      artist.Name,
				Type:       "artist",
			})
		}

		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				results = append(results, SearchResult{
					ArtistID:   artist.ID,
					ArtistName: artist.Name,
					Match:      member,
					Type:       "member",
				})
			}
		}

		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			results = append(results, SearchResult{
				ArtistID:   artist.ID,
				ArtistName: artist.Name,
				Match:      artist.FirstAlbum,
				Type:       "firstAlbum",
			})
		}

		creationDate := fmt.Sprintf("%d", artist.CreationDate)
		if strings.Contains(strings.ToLower(creationDate), query) {
			results = append(results, SearchResult{
				ArtistID:   artist.ID,
				ArtistName: artist.Name,
				Match:      creationDate,
				Type:       "creationDate",
			})
		}
	}

	return results
}
