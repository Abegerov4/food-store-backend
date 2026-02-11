const category = new URLSearchParams(window.location.search).get("name");

fetch(`/api/products?category=${category}`)