const API = "http://localhost:8080";
const params = new URLSearchParams(window.location.search);
const foodId = params.get("id");

function goBack() {
    window.location.href = "index.html";
}

async function updateFood() {
    await fetch(API + "/foods/" + foodId, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            name: foodName.value,
            price: Number(foodPrice.value)
        })
    });

    alert("Food updated");
    goBack();
}