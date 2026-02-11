const API = "http://localhost:8080";

async function register() {
    if (!email.value || !password.value) {
        error.innerText = "All fields are required";
        return;
    }

    const res = await fetch(API + "/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            email: email.value,
            password: password.value
        })
    });

    if (res.ok) {
        window.location.href = "login.html";
    } else {
        error.innerText = "User already exists";
    }
}