document.addEventListener("DOMContentLoaded", loadOrders);

async function loadOrders() {
    const token = localStorage.getItem("token");

    const res = await fetch("/api/admin/orders/all", {
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
    container.innerHTML = "";

    if (!orders || orders.length === 0) {
        container.innerHTML = "<p>No orders yet</p>";
        return;
    }

    orders.forEach(o => {

        const itemsHTML = o.items.map(item => `
            <div class="order-item">
                <span>${item.name}</span>
                <span>₸ ${item.price}</span>
            </div>
        `).join("");

        container.innerHTML += `
            <div class="order-card">
                
                <div class="order-header">
                    <div>
                        <strong>Order ID:</strong>
                        <span>${o.id}</span>
                    </div>

                    <span class="order-status ${o.status}">
                        ${o.status}
                    </span>
                </div>

                <div class="order-body">
                    <p><strong>Email:</strong> ${o.email}</p>
                    <div class="order-items">
                        ${itemsHTML}
                    </div>
                </div>

                <div class="order-footer">
                    <strong>Total: ₸ ${o.total}</strong>

                    ${
                        o.status === "processing"
                        ? `<button class="btn-green"
                            onclick="markCompleted('${o.id}')">
                            Mark completed
                           </button>`
                        : ""
                    }
                </div>

            </div>
        `;
    });
}

async function markCompleted(id) {
    const token = localStorage.getItem("token");

    await fetch(`/api/admin/orders?id=${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({ status: "completed" })
    });

    loadOrders();
}