const category = "bakery";

document.addEventListener("DOMContentLoaded", () => {
    loadProducts();

    const sortSelect = document.getElementById("sortSelect");
    sortSelect.addEventListener("change", () => {
        loadProducts(sortSelect.value);
    });
});

function loadProducts(sort = "", page = 1, limit = 20) {

    let url = `/api/products?category=${category}&page=${page}&limit=${limit}`;

    if (sort) {
        url += `&sort=${sort}`;
    }

    console.log("Fetching:", url);

    fetch(url)
        .then(res => {
            if (!res.ok) {
                throw new Error("Server error: " + res.status);
            }
            return res.json();
        })
        .then(result => {

            console.log("Response:", result);

            let products = [];

            // Старый формат (массив)
            if (Array.isArray(result)) {
                products = result;
            }
            // Новый формат (pagination)
            else if (result && Array.isArray(result.data)) {
                products = result.data;
            }
            else {
                console.error("Unexpected response format:", result);
                products = [];
            }

            renderProducts(products);
        })
        .catch(err => {
            console.error("Fetch error:", err);
            renderProducts([]);
        });
}

function renderProducts(products) {

    if (!Array.isArray(products)) {
        console.error("renderProducts received non-array:", products);
        return;
    }

    const container = document.getElementById("products");
    if (!container) return;

    container.innerHTML = "";

    if (products.length === 0) {
        container.innerHTML = "<p>No products found</p>";
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