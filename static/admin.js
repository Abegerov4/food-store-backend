const API = "http://localhost:8080";
const token = localStorage.getItem("token");

function parseJwt(token) {
	try {
		return JSON.parse(atob(token.split(".")[1]));
	} catch {
		return null;
	}
}

const user = parseJwt(token);

if (!user || user.role !== "admin") {
	window.location.href = "/";
}
// 🔹 LOAD PRODUCTS
async function loadProducts() {
	const res = await fetch(API + "/api/products");
	const data = await res.json();

	const container = document.getElementById("products");
	container.innerHTML = "";

	data.forEach(p => {
		container.innerHTML += `
			<div class="admin-product-card">
				<img src="images/products/${p.image}">

				<div class="admin-product-info">
					<input value="${p.name}" id="name-${p.id}">
					<input value="${p.price}" id="price-${p.id}">
					<input value="${p.category}" id="category-${p.id}">
				</div>

				<div class="admin-actions">
					<button class="btn-blue btn-small"
						onclick="updateProduct('${p.id}')">
						Update
					</button>

					<button class="btn-red btn-small"
						onclick="deleteProduct('${p.id}')">
						Delete
					</button>
				</div>
			</div>
		`;
	});
}

// ➕ ADD
async function addProduct() {
	const res = await fetch(API + "/api/admin/products", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			"Authorization": "Bearer " + token
		},
		body: JSON.stringify({
			name: productName.value,
			price: Number(productPrice.value),
			category: productCategory.value,
			image: productImage.value
		})
	});

	if (res.ok) {
		loadProducts();
		alert("Product added ✅");
	} else {
		alert(await res.text());
	}
}
//update
async function updateProduct(id) {
    await fetch(API + "/api/admin/products/" + id, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({
            name: document.getElementById(`name-${id}`).value,
            price: Number(document.getElementById(`price-${id}`).value),
            category: document.getElementById(`category-${id}`).value
        })
    });

    alert("Updated ✅");
    loadProducts();
}

// ❌ DELETE
async function deleteProduct(id) {
	if (!confirm("Delete product?")) return;

	await fetch(API + "/api/admin/products/" + id, {
		method: "DELETE",
		headers: {
			"Authorization": "Bearer " + token
		}
	});

	loadProducts();
}

document.addEventListener("DOMContentLoaded", loadProducts);