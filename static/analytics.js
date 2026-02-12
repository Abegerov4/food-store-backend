document.addEventListener("DOMContentLoaded", () => {

    const token = localStorage.getItem("token");

    if (!token) {
        window.location.href = "/";
        return;
    }

    const payload = JSON.parse(atob(token.split(".")[1]));

    if (payload.role !== "admin") {
        window.location.href = "/";
        return;
    }

    loadAnalytics();
});

async function loadAnalytics() {

    const token = localStorage.getItem("token");

    const res = await fetch("/api/admin/analytics", {
        headers: {
            "Authorization": "Bearer " + token
        }
    });

    if (!res.ok) {
        alert("Unauthorized");
        return;
    }

    const data = await res.json();

    // Revenue
    document.getElementById("revenue").innerText =
        data.totalRevenue + " ₸";

    // Orders
    document.getElementById("orders").innerText =
        data.totalOrders;

    // Top products
    const list = document.getElementById("topProducts");
    list.innerHTML = "";

    data.topProducts.forEach(p => {
        list.innerHTML += `
            <li>
                ${p._id} — ${p.quantity} sales
            </li>
        `;
    });
}
const token = localStorage.getItem("token");

fetch("/api/admin/analytics/revenue", {
    headers: {
        "Authorization": "Bearer " + token
    }
})
.then(res => res.json())
.then(data => {

    const labels = data.map(d => d._id);
    const revenues = data.map(d => d.revenue);

    const ctx = document.getElementById("revenueChart").getContext("2d");

    new Chart(ctx, {
        type: "line",
        data: {
            labels: labels,
            datasets: [{
                label: "Revenue ₸",
                data: revenues,
                borderColor: "#10b981",
                backgroundColor: "rgba(16,185,129,0.1)",
                tension: 0.4,
                fill: true
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: { display: true }
            },
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    });

});