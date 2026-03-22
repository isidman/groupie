# Groupie Tracker

A web application built in Go that consumes the [Groupie Trackers API](https://groupietrackers.herokuapp.com/api) and displays artist and concert information through a clean, concert poster-inspired interface.

## Stack

- **Backend:** Go (standard library only)
- **Frontend:** HTML, CSS, JavaScript (no frameworks)
- **Database:** None — data is fetched from the external API and cached in memory
- **Fonts:** Bebas Neue, DM Mono (Google Fonts)
- **Map:** Leaflet.js (client-side, loaded from CDN)

## Running the project

```bash
go run .
```

Then open `http://localhost:8080` in your browser.

---

## Project Structure

```
groupie-tracker/
├── internal/           # Backend logic (API client, cache, models, search, filter)
│   ├── cache.go
│   ├── client.go
│   ├── filter.go
│   ├── merge.go
│   ├── models.go
│   └── search.go
├── web/                # HTTP handlers and templates
│   ├── errors.go
│   ├── handlers.go
│   ├── middleware.go
│   └── templates.go
├── static/             # Frontend assets
│   ├── app.js
│   └── style.css
├── web/                # HTML templates
│   ├── index.html
│   ├── artist.html
│   └── error.html
└── main.go
```

---

## Core Features

- Displays all artists from the API in a responsive card grid
- Individual artist pages showing members, creation year, first album, concert dates and locations
- Client-server event system: a button on each artist page fetches and displays the full artist JSON data asynchronously via `fetch()`
- Error handling for 404 and 500 HTTP status codes
- Server never crashes — all panics are recovered by middleware

---

## Optional: Visualizations

The entire frontend follows a **concert poster aesthetic**:

- Dark background (`#1a1a18`) with yellow accent (`#FFED72`)
- Bebas Neue for headings, DM Mono for body text
- Cards lift on hover with a yellow border accent
- Concert blocks styled with a yellow left stripe like ticket stubs
- Fully responsive — artist detail page stacks vertically on mobile via a CSS media query breakpoint at 600px
- Styled 404 and 500 error pages consistent with the overall design
- Good color contrast throughout (light text on dark backgrounds, yellow on dark for accents)

---

## Optional: Search Bar

A live search bar accessible from the homepage header.

**How it works:**
- Click the 🔍 button — the search bar slides into view
- Type at least 2 characters to trigger suggestions
- Press `Escape` or click outside to close

**Searches across:**
- Artist / band name
- Member names
- Concert locations (raw format, e.g. `saitama-japan`)
- First album date
- Creation date

**Features:**
- Case-insensitive matching
- Each suggestion shows the matched text and its type (Artist, Member, Location, First Album, Creation Date)
- Clicking a suggestion navigates directly to the artist's page
- Location data is pre-fetched and cached at server startup for fast search

**API endpoint:** `GET /api/search?q=<query>`

---

## Optional: Filters

A filter modal accessible from the homepage header via the ⚡ button.

**How it works:**
- Click ⚡ to open the filter panel overlay
- Adjust filters — results update live without page reload (asynchronous)
- Press `Escape` or click outside the panel to close

**Filter types:**
- **Creation year range** — range input (e.g. 1995 to 2000)
- **First album year range** — range input (e.g. 1990 to 1992)
- **Number of members** — checkboxes (1, 2, 3, 4, 5, 6, 7+)
- **Concert location** — text input with substring matching

**Results** appear as a scrollable list inside the modal with artist thumbnail, name, and creation year. Clicking a result navigates to the artist's page.

**API endpoint:** `GET /api/filter?creationMin=&creationMax=&albumMin=&albumMax=&members=&location=`

---

## Optional: Geolocalization

Each artist page includes an interactive map showing their concert locations.

**How it works:**
- Raw location strings (e.g. `nagoya-japan`, `abu_dhabi-united_arab_emirates`) are stored alongside artist data in the cache
- When an artist page loads, the JS fetches coordinates for each location from the `/api/geocode` endpoint
- Requests are rate-limited to one per 1.1 seconds to comply with Nominatim's usage policy
- Markers appear on the map sequentially as coordinates are resolved
- Clicking a marker shows a popup with the location name

**Geocoding logic:**
- Single-word country suffix (e.g. `japan`, `usa`) → country placed first in query for better Nominatim results
- Multi-word country suffix (e.g. `united_arab_emirates`) → forward order used

**Libraries used:**
- [Nominatim](https://nominatim.openstreetmap.org) — free geocoding API, called from Go, no API key required
- [Leaflet.js](https://leafletjs.com) — open-source JS mapping library, loaded from CDN
- OpenStreetMap tiles

**API endpoint:** `GET /api/geocode?location=<raw-location-string>`

---

## Authors

 [G.Argyros](https://github.com/GeorgeArgygos)

 [I.Manolas](https://github.com/isidman)