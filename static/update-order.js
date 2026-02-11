const API = "http://localhost:8080";

const params = new URLSearchParams(window.location.search);
const orderId = params.get("id");

function goBack() {
    window.location.href = "index.html";
}

async function updateOrder() {
    await fetch(API + "/orders/" + orderId, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            food_ids: orderFoods.value.split(",").map(x => x.trim())
        })
    });

    alert("Order updated successfully");
    window.location.href = "index.html";
}