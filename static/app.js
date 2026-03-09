/*
app.js
JavaScript πλευράς client για δυναμικά events/actions.
Χρησιμοποιείται για επικοινωνία client-server χωρίς ανανέωση σελίδας.
*/
document.addEventListener("DOMContentLoaded", function () {
  const button = document.getElementById("load-data-btn");
  const output = document.getElementById("artist-json-output");

  if (!button || !output) {
    return;
  }

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
});