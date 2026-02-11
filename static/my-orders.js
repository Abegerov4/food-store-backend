document.addEventListener("DOMContentLoaded", loadOrders);

async function loadOrders() {
    const token = localStorage.getItem("token");

    if (!token) {
        alert("Please login first");
        return;
    }

    const res = await fetch("/api/orders", {
        headers: {
            "Authorization": "Bearer " + token
        }
    });

    if (!res.ok) {
        alert("Failed to load orders");
        return;
    }

    const orders = await res.json();

    const container = document.getElementById("orders");

    if (!container) return;

    container.innerHTML = "";

    // 🔒 защита от null
    if (!orders || orders.length === 0) {
        container.innerHTML = "<p>No orders yet</p>";
        return;
    }

    orders.forEach(o => {
        container.innerHTML += `
            <div class="cart-item">
                <div class="cart-info">
                    <h3>Order</h3>
                    <p>Total: ₸ ${o.total}</p>
                    <p>Status: ${o.status}</p>
                </div>
            </div>
        `;
    });
}