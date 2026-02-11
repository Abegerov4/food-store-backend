const category = "pets";

document.addEventListener("DOMContentLoaded", () => {
    loadProducts();

    const sortSelect = document.getElementById("sortSelect");
    sortSelect.addEventListener("change", () => {
        loadProducts(sortSelect.value);
    });
});

function loadProducts(sort = "") {

    let url = `/api/products?category=${category}`;

    if (sort) {
        url += `&sort=${sort}`;
    }

    console.log("Fetching:", url);

    fetch(url)
        .then(res => res.json())
        .then(renderProducts)
        .catch(console.error);
}

function renderProducts(products) {
    const container = document.getElementById("products");
    container.innerHTML = "";

    if (products.length === 0) {
        container.innerHTML = "<p>No pet products yet</p>";
        return;
    }

    products.forEach(p => {
        container.innerHTML += `
            <div class="product-card">
                <img src="../images/products/${p.image}">
                <h3>${p.name}</h3>
                <p class="price">₸ ${p.price}</p>
                <button class="btn-primary"
                    onclick='addToCart(${JSON.stringify(p)})'>
                    Add to cart
                </button>
            </div>
        `;
    });
}

function addToCart(product) {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];
    cart.push(product);
    localStorage.setItem("cart", JSON.stringify(cart));
    alert("Added to cart");
}