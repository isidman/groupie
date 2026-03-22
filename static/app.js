/*
app.js
JavaScript πλευράς client για δυναμικά events/actions.
Χρησιμοποιείται για επικοινωνία client-server χωρίς ανανέωση σελίδας.
*/
document.addEventListener("DOMContentLoaded", function () {

  // Load JSON button (artist page only)
  const button = document.getElementById("load-data-btn")
  const output = document.getElementById("artist-json-output")

  if (button && output) {
    button.addEventListener("click", async function () {
      const artistID = button.getAttribute("data-artist-id")
      output.textContent = "Loading..."
      try {
        const response = await fetch(`/api/artist?id=${artistID}`)
        if (!response.ok) {
          output.textContent = `Request failed with status ${response.status}`
          return
        }
        const data = await response.json()
        output.textContent = JSON.stringify(data, null, 2)
      } catch (error) {
        output.textContent = "An error occurred while loading artist data."
      }
    })
  }

  // Search toggle (show/hide search bar)
  const searchInput = document.getElementById("search-input")
  const suggestions = document.getElementById("search-suggestions")
  const searchBtn = document.getElementById("search-btn")
  const searchWrap = document.getElementById("search-wrap")

  if (searchBtn && searchWrap) {
    searchBtn.addEventListener("click", function () {
      searchBtn.classList.add("hidden")
      searchWrap.classList.add("open")
      searchInput.focus()
    })
  }

  // Search suggestions
  if (searchInput && suggestions) {

    // Fetch and render suggestions as user types
    searchInput.addEventListener("input", async function () {
      const query = searchInput.value.trim()

      // If query is less than 2 characters, hide suggestions and return
      if (query.length < 2) {
        suggestions.innerHTML = ""
        suggestions.classList.remove("open")
        return
      }

      // Fetch from the search endpoint
      const response = await fetch(`/api/search?q=${query}`)
      const results = await response.json()

      if (results.length === 0) {
        suggestions.innerHTML = ""
        suggestions.classList.remove("open")
        return
      }

      // Build the dropdown HTML
      suggestions.innerHTML = ""
      for (const result of results) {
        const item = document.createElement("div")
        item.className = "suggestion-item"
        item.innerHTML = `
          <span>${result.artistName} - ${result.match}</span>
          <span class="suggestion-type">${result.type}</span>
        `
        // Clicking navigates to artist page
        item.addEventListener("click", function () {
          window.location.href = `/artist?id=${result.artistID}`
        })
        suggestions.appendChild(item)
      }

      suggestions.classList.add("open")
    })

    // Esc closes the search bar
    searchInput.addEventListener("keydown", function (e) {
      if (e.key === "Escape") {
        searchWrap.classList.remove("open")
        searchBtn.classList.remove("hidden")
        searchInput.value = ""
        suggestions.innerHTML = ""
        suggestions.classList.remove("open")
      }
    })

    // Close suggestions when clicking outside
    document.addEventListener("click", function (e) {
      if (!searchInput.contains(e.target) && !suggestions.contains(e.target)) {
        suggestions.classList.remove("open")
      }

      // Close search bar when clicking outside the whole search row
      if (searchWrap && searchBtn) {
        const searchRow = document.querySelector(".search-row")
        if (searchRow && !searchRow.contains(e.target)) {
          searchWrap.classList.remove("open")
          searchBtn.classList.remove("hidden")
          searchInput.value = ""
          suggestions.innerHTML = ""
          suggestions.classList.remove("open")
        }
      }
    })
  }

  // Filter modal
  const filterBtn = document.getElementById("filter-btn")
  const filterModal = document.getElementById("filter-modal")
  const filterClose = document.getElementById("filter-close")
  const filterResults = document.getElementById("filter-results")
  const filterLocation = document.getElementById("filter-location")

  if (filterBtn && filterModal) {

    // Open
    filterBtn.addEventListener("click", function () {
      filterModal.classList.add("open")
    })

    // Close
    filterClose.addEventListener("click", function () {
      filterModal.classList.remove("open")
    })

    // Clicking outside the panel closes it
    filterModal.addEventListener("click", function (e) {
      if (e.target === filterModal) {
        filterModal.classList.remove("open")
      }
    })

    // Esc closes the panel
    document.addEventListener("keydown", function(e) {
      if (e.key === "Escape") {
        filterModal.classList.remove("open")
      }
    })

    // Filters run on change of any input
    async function runFilter() {
      const params = new URLSearchParams()

      const creationMin = document.getElementById("creationMin").value
      const creationMax = document.getElementById("creationMax").value
      const albumMin = document.getElementById("albumMin").value
      const albumMax = document.getElementById("albumMax").value

      if (creationMin) { params.append("creationMin", creationMin) }
      if (creationMax) { params.append("creationMax", creationMax) }
      if (albumMin) { params.append("albumMin", albumMin) }
      if (albumMax) { params.append("albumMax", albumMax) }

      const checkedMembers = document.querySelectorAll("#members-checkboxes input:checked")
      checkedMembers.forEach(function (cb) {
        if (cb.value === "7") {
          for (let i = 7; i <= 20; i++) {
            params.append("members", i)
          }
        } else {
          params.append("members", cb.value)
        }
      })

      const location = filterLocation.value.trim()
      if (location) { params.append("location", location) }

      // Fetch results
      const response = await fetch(`/api/filter?${params.toString()}`)
      const artists = await response.json()

      // Render them
      filterResults.innerHTML = ""

      if (artists.length === 0) {
        filterResults.innerHTML = `<p class="filter-hint">No artists match these filters.</p>`
        return
      }

      for (const artist of artists) {
        const item = document.createElement("a")
        item.className = "filter-result-item"
        item.href = `/artist?id=${artist.id}`
        item.innerHTML = `
          <img class="filter-result-thumb" src="${artist.image}" alt="${artist.name}" />
          <div class="filter-result-info">
            <strong>${artist.name}</strong>
            <p>Est. ${artist.creationDate} · ${artist.members.length} members</p>
          </div>
        `
        filterResults.appendChild(item)
      }
    }

    //runFilter is attached to all inputs
    document.getElementById("creationMin").addEventListener("input", runFilter)
    document.getElementById("creationMax").addEventListener("input", runFilter)
    document.getElementById("albumMin").addEventListener("input", runFilter)
    document.getElementById("albumMax").addEventListener("input", runFilter)
    filterLocation.addEventListener("input", runFilter)

    document.querySelectorAll("#members-checkboxes input").forEach(function (cb) {
      cb.addEventListener("change", runFilter)
    })

  }

  const mapEl= document.getElementById("concert-map")

    if (mapEl) {
      // Read locations from data attribute
      const raw = mapEl.getAttribute("data-locations") || ""
      const locations = raw.split("|").map(l => l.trim()).filter(l => l.length > 0)

      // Start leaflet map
      const map = L.map("concert-map").setView([20, 0], 2)
      
      //OSM layer
      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution: "© OpenStreetMap contributors"
      }).addTo(map)

      // Geocode locations and add markers
      //Rate-limited by Nominatim: 1 request per second
      locations.forEach(function(location, index) {
        setTimeout(async function() {
          try {
            const response = await fetch(`/api/geocode?location=${encodeURIComponent(location)}`)
            if (!response.ok) return
            const geo = await response.json()
            if (geo.lat && geo.lon) {
              const lat = parseFloat(geo.lat)
              const lon = parseFloat(geo.lon)
              L.marker([lat,lon])
                .addTo(map)
                .bindPopup(`<b>${location}</b>`)
            }
          } catch (e) {
            console.error("Geocode failed for", location, e)
          }
        }, index *1100) // 1.1s Delay between requests
      })
    }
})