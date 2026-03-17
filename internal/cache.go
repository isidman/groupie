package internal

import (
	"log"
	"sync"
	"time"
)

type ArtistCache struct {
	mu      sync.RWMutex
	artists []Artist
	client  *Client
	refresh chan struct{}
}

func NewArtistCache(client *Client) *ArtistCache {
	return &ArtistCache{
		client:  client,
		refresh: make(chan struct{}, 1),
	}
}

func (c *ArtistCache) Start() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		c.load()

		for {
			select {
			case <-ticker.C:
				c.load()
			case <-c.refresh:
				c.load()
			}
		}
	}()
}

func (c *ArtistCache) load() {
	artists, err := c.client.GetArtists()
	if err != nil {
		log.Println("cache refresh failed:", err)
		return
	}

	for i, artist := range artists {
		locations, err := c.client.GetLocationsByID(artist.ID)
		if err != nil {
			log.Println("location fetch failed for artist", artist.ID, err)
			continue //skip it, don't crash
		}
		artists[i].ConcertLocations = locations.Locations
	}

	c.mu.Lock()
	c.artists = artists
	c.mu.Unlock()

}

func (c *ArtistCache) GetArtists() []Artist {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]Artist, len(c.artists))
	copy(result, c.artists)
	return result
}

func (c *ArtistCache) RequestRefresh() {
	select {
	case c.refresh <- struct{}{}:
	default:
	}
}
