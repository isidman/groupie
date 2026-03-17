/*
app.js
JavaScript πλευράς client για δυναμικά events/actions.
Χρησιμοποιείται για επικοινωνία client-server χωρίς ανανέωση σελίδας.
*/
document.addEventListener("DOMContentLoaded", function () {
  const button = document.getElementById("load-data-btn");
  const output = document.getElementById("artist-json-output");


  if (button && output) {
      button.addEventListener("click", async function () {
      const artistID = button.getAttribute("data-artist-id");

      output.textContent = "Loading...";

      try {
        const response = await fetch(`/api/artist?id=${artistID}`);

        if (!response.ok) {
          output.textContent = `Request failed with status ${response.status}`;
          return;
        }

        const data = await response.json();
       output.textContent = JSON.stringify(data, null, 2);
      } catch (error) {
        output.textContent = "An error occurred while loading artist data.";
      }
    });
  }
  
  
  
  const searchInput = document.getElementById("search-input");
  const suggestions = document.getElementById("search-suggestions");

  if (searchInput && suggestions) {
  //if they don't exist, return early
  
  searchInput.addEventListener("input", async function () {
    const query = searchInput.value.trim()

    //if query is less than 2 characters, hide suggestions and return
    if (query.length < 2) {
      suggestions.innerHTML = ""
      suggestions.classList.remove("open")
      return
    }

    //fetch from the search endpoint
    const response = await fetch(`/api/search?q=${query}`)
    const results = await response.json()

    if (results.length === 0) {
      suggestions.innerHTML = ""
      suggestions.classList.remove("open")
      return
    }

    //build the dropdown HTML
    suggestions.innerHTML = ""
    for (const result of results) {
      const item = document.createElement("div")
      item.className = "suggestion-item"
      item.innerHTML = `
        <span>${result.artistName} - ${result.match}</span>
        <span class="suggestion-type">${result.type}</span>
        `

        //clicking navigates to artist page
        item.addEventListener("click", function(){
          window.location.href = `/artist?id=${result.artistID}`
        })
        suggestions.appendChild(item)
    }

    suggestions.classList.add("open")
  })

    //close suggestions when clicking outside
    document.addEventListener("click", function(e){
      if (!searchInput.contains(e.target) && !suggestions.contains(e.target)) {
        suggestions.classList.remove("open")
      }
    })
  }
});