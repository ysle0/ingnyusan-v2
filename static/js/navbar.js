// Navbar scroll behavior using IntersectionObserver
document.addEventListener('DOMContentLoaded', function() {
  const navbar = document.querySelector('header nav');
  
  // Create a sentinel element for observation
  const sentinel = document.createElement('div');
  sentinel.setAttribute('id', 'navbar-sentinel');
  
  // Style it as invisible but positioned at 15% down the page
  sentinel.style.position = 'absolute';
  sentinel.style.top = '15vh'; // 15% of viewport height
  sentinel.style.left = '0';
  sentinel.style.width = '100%';
  sentinel.style.height = '1px';
  sentinel.style.pointerEvents = 'none';
  sentinel.style.opacity = '0';
  
  // Insert it at the beginning of the body
  document.body.insertBefore(sentinel, document.body.firstChild);
  
  // Set up the IntersectionObserver
  const observer = new IntersectionObserver(
    (entries) => {
      // When sentinel goes out of view (scrolled past 15%)
      if (!entries[0].isIntersecting) {
        navbar.classList.add('scrolled');
      } else {
        navbar.classList.remove('scrolled');
      }
    },
    {
      // Observer options
      rootMargin: '0px',
      threshold: 0 // Trigger as soon as any part is invisible
    }
  );
  
  // Start observing the sentinel element
  observer.observe(sentinel);
  
  // Check initial position in case page loads scrolled down
  if (window.scrollY > window.innerHeight * 0.15) {
    navbar.classList.add('scrolled');
  }
}); 