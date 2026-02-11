const API = "http://localhost:8080";

async function addFood() {
    await fetch(API + "/foods", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            id: foodId.value,
            name: foodName.value,
            price: Number(foodPrice.value)
        })
    });

    alert("Food added. Click 'Load Foods' to see it.");
}

// LOAD FOODS
async function loadFoods() {
    const res = await fetch(API + "/foods");
    const data = await res.json();

    foods.innerHTML = "";
    data.forEach(f => {
        foods.innerHTML += `
            <div class="item">
                <span>${f.name} — ${f.price}₸</span>

                <div class="actions">
                    <button class="btn-small btn-blue"
                        onclick="openUpdateFood('${f.id}')">
                        Update
                    </button>

                    <button class="btn-small btn-red"
                        onclick="deleteFood('${f.id}')">
                        Delete
                    </button>
                </div>
            </div>
        `;
    });
}

// OPEN UPDATE FOOD PAGE
function openUpdateFood(id) {
    window.location.href = `update-food.html?id=${id}`;
}

// DELETE FOOD
async function deleteFood(id) {
    await fetch(API + "/foods/" + id, {
        method: "DELETE"
    });

    loadFoods();
}


async function createOrder() {
    await fetch(API + "/orders", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            id: orderId.value,
            food_ids: orderFoods.value
                .split(",")
                .map(x => x.trim())
        })
    });

    alert("Order created. Click 'Load Orders' to see status.");
}

// LOAD ORDERS
async function loadOrders() {
    const res = await fetch(API + "/orders");
    const data = await res.json();

    orders.innerHTML = "";
    data.forEach(o => {
        orders.innerHTML += `
            <div class="item">
                <div class="item-left">
                <span>Order ${o.id} — ${o.total}₸</span>
                <span class="badge ${o.status}">${o.status}</span>
                </div>

                <div class="actions">
                <button class="btn-small btn-blue"
                    onclick="openUpdateOrder('${o.id}')">
                    Update
                </button>

                <button class="btn-small btn-red"
                    onclick="deleteOrder('${o.id}')">
                    Delete
                </button>
                </div>
            </div>
            `;  
    });
}

// OPEN UPDATE ORDER PAGE
function openUpdateOrder(id) {
    window.location.href = `update-order.html?id=${id}`;
}

// DELETE ORDER
async function deleteOrder(id) {
    await fetch(API + "/orders/" + id, {
        method: "DELETE"
    });

    loadOrders();
}
function goLogin() {
    window.location.href = "login.html";
}

function goRegister() {
    window.location.href = "register.html";
}
document.addEventListener("DOMContentLoaded", updateCartCount);

function updateCartCount() {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];
    const countEl = document.getElementById("cartCount");
    if (countEl) {
        countEl.innerText = cart.length;
    }
}
function parseJwt(token) {
    try {
        const base64 = token.split(".")[1];
        return JSON.parse(atob(base64));
    } catch {
        return null;
    }
}
function renderAuth() {
    const authArea = document.getElementById("authArea");
    const token = localStorage.getItem("token");

    if (!token) {
        // ❌ НЕ залогинен
        authArea.innerHTML += `
            <a href="/login.html" class="link-btn">Login</a>
            <a href="/register.html" class="link-btn">Register</a>
        `;
        return;
    }

    const user = parseJwt(token);
    if (!user) {
        localStorage.removeItem("token");
        location.reload();
        return;
    }

    // ✅ ЗАЛОГИНЕН
    authArea.innerHTML += `
        <span class="user-email">${user.email}</span>
        <button class="link-btn" onclick="logout()">Logout</button>
    `;
}
function logout() {
    localStorage.removeItem("token");
    window.location.href = "/index.html";
}
document.addEventListener("DOMContentLoaded", () => {
    updateCartCount();
    renderAuth();
});
function searchFoods() {
    // TODO: реализовать поиск
}

// 🔥 POPULAR PRODUCTS
document.addEventListener("DOMContentLoaded", loadPopularProducts);

async function loadPopularProducts() {
    const res = await fetch("/api/products");
    const products = await res.json();

    const container = document.getElementById("popularProducts");
    if (!container) return;

    container.innerHTML = "";

    // берем первые 15
    products.slice(0, 15).forEach(p => {
        container.innerHTML += `
            <div class="product-card">
                <img src="images/products/${p.image}">
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
async function searchProducts() {
    const query = document.getElementById("searchInput").value;

    if (!query) return;

    const res = await fetch(`/api/products?search=${query}`);
    const data = await res.json();

    renderSearchResults(data);
}
function renderSearchResults(products) {
    const section = document.getElementById("searchResultsSection");
    const container = document.getElementById("searchResults");

    section.style.display = "block";
    container.innerHTML = "";

    if (products.length === 0) {
        container.innerHTML = "<p>No products found</p>";
        return;
    }

    products.forEach(p => {
        container.innerHTML += `
            <div class="product-card">
                <img src="images/products/${p.image}">
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
let searchTimeout = null;

function liveSearch() {
    const query = document.getElementById("searchInput").value.trim();

    // очищаем предыдущий таймер
    clearTimeout(searchTimeout);

    // если пусто → скрыть результаты
    if (!query) {
        document.getElementById("searchResultsSection").style.display = "none";
        return;
    }

    // задержка 300мс (чтобы не спамить сервер)
    searchTimeout = setTimeout(() => {
        fetch(`/api/products?search=${query}`)
            .then(res => res.json())
            .then(renderSearchResults);
    }, 300);
}
function renderAuth() {
    const authArea = document.getElementById("authArea");
    const token = localStorage.getItem("token");

    let html = `
        <a href="/cart.html" class="cart-icon">
            🛒
            <span id="cartCount" class="cart-count">0</span>
        </a>
    `;

    if (!token) {
        html += `
            <a href="/login.html" class="link-btn">Login</a>
            <a href="/register.html" class="link-btn">Register</a>
        `;
    } else {
        const user = parseJwt(token);

        if (!user) {
            localStorage.removeItem("token");
            location.reload();
            return;
        }

        html += `
            <span class="user-email">${user.email}</span>
            <button class="link-btn" onclick="logout()">Logout</button>
        `;

        if (user.role === "admin") {
            html += `
                <a href="/admin.html" class="link-btn">Admin</a>
            `;
        }
    }

    authArea.innerHTML = html;
    updateCartCount();
}
document.addEventListener("DOMContentLoaded", () => {
    renderAuth();
});
function addToCart(product) {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];
    cart.push(product);
    localStorage.setItem("cart", JSON.stringify(cart));
    updateCartCount();
    alert("Added to cart");
}