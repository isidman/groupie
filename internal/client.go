// models.go
// Περιέχει όλα τα structs που αντιστοιχούν
// στα JSON responses του εξωτερικού Groupie Trackers API.
//
// Τα structs αυτά χρησιμοποιούνται μόνο για αποθήκευση δεδομένων
// και όχι για απευθείας εμφάνιση στα templates.

package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseAPI = "https://groupietrackers.herokuapp.com/api"

type Client struct {
	httpClient *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
	}
}

// fetchJSON performs a GET request and decodes JSON into v
func (c *Client) fetchJSON(url string, v any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status from API: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}

func (c *Client) GetArtists() ([]Artist, error) {
	var artists []Artist
	if err := c.fetchJSON(baseAPI+"/artists", &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func (c *Client) GetRelationByID(id int) (Relation, error) {
	var rel Relation
	if err := c.fetchJSON(fmt.Sprintf("%s/relation/%d", baseAPI, id), &rel); err != nil {
		return Relation{}, err
	}
	return rel, nil
}

func (c *Client) GetLocationsByID(id int) (Locations, error) {
	var loc Locations
	if err := c.fetchJSON(fmt.Sprintf("%s/locations/%d", baseAPI, id), &loc); err != nil {
		return Locations{}, err
	}
	return loc, nil
}

func (c *Client) GetDatesByID(id int) (Dates, error) {
	var dates Dates
	if err := c.fetchJSON(fmt.Sprintf("%s/dates/%d", baseAPI, id), &dates); err != nil {
		return Dates{}, err
	}
	return dates, nil
}
