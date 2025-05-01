
// Theme toggle functionality
document.addEventListener('DOMContentLoaded', function () {
  const sunIcon = document.getElementById('sun-icon');
  const moonIcon = document.getElementById('moon-icon');
  const themeToggle = document.getElementById('theme-toggle');

  // Check for saved theme preference in cookie or default
  const currentTheme = getCookie("theme") || 'light';
  // if (currentTheme) {
  // console.log("Current theme:", currentTheme); // Will show "light" or "dark"
  // }

  // Set initial icon state
  updateIconVisibility(currentTheme);

  // Set the body theme
  document.documentElement.setAttribute('data-theme', currentTheme);
  // document.body.setAttribute('data-theme', currentTheme);

  // Log initial state for debugging
  // console.log('Initial theme:', currentTheme);
  // console.log('HTML data-theme attribute:', document.documentElement.getAttribute('data-theme'));
  //

  // Function to update icon visibility
  function updateIconVisibility(theme) {
    if (theme === 'light') {
      if (sunIcon) sunIcon.style.display = 'inline';
      if (moonIcon) moonIcon.style.display = 'none';
    } else {
      if (sunIcon) sunIcon.style.display = 'none';
      if (moonIcon) moonIcon.style.display = 'inline';
    }
  }


  // Add event listener to toggle button
  themeToggle.addEventListener('click', function (e) {
    e.preventDefault();

    // Toggle theme directly first for immediate feedback
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'light';
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    setThemeWatercss(newTheme);

    // Apply the new theme immediately
    document.documentElement.setAttribute('data-theme', newTheme);
    updateIconVisibility(newTheme);

    // console.log('Theme toggled to:', newTheme);

    // Then call the API to save the preference
    fetch('/api/toggle-theme')
      .then(response => response.json())
      .then(data => {
        // console.log('API response:', data);
        // The theme is already updated for better UX
      })
      .catch(err => {
        console.error('Error toggling theme:', err);
      });
  });
});
// Function to get a specific cookie by name
function getCookie(name) {
  // Add "=" to the name to make sure we match the exact cookie
  const cookieName = name + "=";
  // Get all cookies in document.cookie (returns a string of all cookies, separated by semicolons)
  const cookies = document.cookie.split(';');

  // Loop through all cookies
  for (let i = 0; i < cookies.length; i++) {
    // Get a single cookie string
    let cookie = cookies[i].trim();
    // Check if this cookie starts with the name we're looking for
    if (cookie.indexOf(cookieName) === 0) {
      // Return the cookie value (everything after the equals sign)
      return cookie.substring(cookieName.length, cookie.length);
    }
  }
  // Return null if cookie not found
  return null;
}

function setThemeWatercss(themeName) {
  const link = document.querySelector('head link');
  link.href = link.href.replace(/\min.css$/, `-${themeName}.min.css`);
}