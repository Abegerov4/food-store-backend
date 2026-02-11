const API = "http://localhost:8080";

async function login() {
    if (!email.value || !password.value) {
        error.innerText = "All fields are required";
        return;
    }

    const res = await fetch(API + "/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            email: email.value,
            password: password.value
        })
    });

    const data = await res.json();

    if (res.ok) {
        // ✅ сохраняем JWT
        localStorage.setItem("token", data.token);
        window.location.href = "index.html";
    } else {
        error.innerText = data.error || "Invalid email or password";
    }
}