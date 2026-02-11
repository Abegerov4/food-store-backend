const category = "fruit"; // 👈 ВАЖНО

document.addEventListener("DOMContentLoaded", loadProducts);

function loadProducts() {
    fetch(`/api/products?category=${category}`)
        .then(res => res.json())
        .then(renderProducts)
        .catch(console.error);
}

function renderProducts(products) {
    const container = document.getElementById("products");
    container.innerHTML = "";

    if (products.length === 0) {
        container.innerHTML = "<p>No products yet</p>";
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