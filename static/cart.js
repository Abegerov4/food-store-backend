document.addEventListener("DOMContentLoaded", loadCart);

function loadCart() {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];
    const container = document.getElementById("cart");
    const totalEl = document.getElementById("total");

    container.innerHTML = "";
    let total = 0;

    if (cart.length === 0) {
        container.innerHTML = "<p>Your cart is empty.</p>";
        totalEl.innerText = "₸ 0";
        return;
    }

    cart.forEach((item, index) => {
        total += item.price;

        container.innerHTML += `
            <div class="cart-item">
                <img src="images/products/${item.image}">
                <div class="cart-info">
                    <h3>${item.name}</h3>
                    <p>₸ ${item.price}</p>
                </div>
                <button class="btn-remove"
                    onclick="removeItem(${index})">
                    Remove
                </button>
            </div>
        `;
    });

    totalEl.innerText = "₸ " + total;
}

function removeItem(index) {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];
    cart.splice(index, 1);
    localStorage.setItem("cart", JSON.stringify(cart));
    loadCart();
}

function clearCart() {
    localStorage.removeItem("cart");
    loadCart();
}

function checkout() {
    const cart = JSON.parse(localStorage.getItem("cart")) || [];
    const token = localStorage.getItem("token");

    if (!token) {
        alert("Please login first");
        window.location.href = "/login.html";
        return;
    }

    if (cart.length === 0) {
        alert("Cart is empty");
        return;
    }

    const order = {
        items: cart,
        total: cart.reduce((sum, p) => sum + p.price, 0)
    };

    fetch("/api/orders", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token   // ✅ ВОТ ЭТО ВАЖНО
        },
        body: JSON.stringify(order)
    })
    .then(res => {
        if (!res.ok) throw new Error("Order failed");
        return res.json();
    })
    .then(() => {
        alert("Order created successfully!");
        localStorage.removeItem("cart");
        window.location.href = "/my-orders.html";
    })
    .catch(() => alert("Checkout failed"));
}