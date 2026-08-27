package internal

import (
	"strconv"
	"strings"
)

type FilterParams struct {
	CreationMin int
	CreationMax int
	AlbumMin    int // Year-part only
	AlbumMax    int
	Members     []int  // List of member counts to include
	Location    string // Location substring to match
}

func FilterArtists(artists []Artist, params FilterParams) []Artist {
	results := []Artist{}

	for _, artist := range artists {
		// Creation date range
		if params.CreationMin > 0 && artist.CreationDate < params.CreationMin {
			continue
		}
		if params.CreationMax > 0 && artist.CreationDate > params.CreationMax {
			continue
		}

		//First album, year range
		//First album from artist array is just a string through A T O I
		albumYearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
		albumYear, _ := strconv.Atoi(albumYearStr)

		if params.AlbumMin > 0 && albumYear < params.AlbumMin {
			continue
		}
		if params.AlbumMax > 0 && albumYear > params.AlbumMax {
			continue
		}

		// Members filter
		if len(params.Members) > 0 {
			found := false
			for _, m := range params.Members {
				if len(artist.Members) == m {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if params.Location != "" {
			found := false
			for _, loc := range artist.ConcertLocations {
				formattedLoc := FormatLocation(loc)
				if strings.Contains(strings.ToLower(loc), strings.ToLower(params.Location)) ||
					strings.Contains(strings.ToLower(formattedLoc), strings.ToLower(params.Location)) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		results = append(results, artist)
	}

	return results
}
