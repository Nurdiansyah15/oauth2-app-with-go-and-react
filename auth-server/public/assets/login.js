document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("myForm");
  const responseElement = document.getElementById("response");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    // Get form data
    const formData = new FormData(form);
    const formObject = {};
    const clientId = formData.get("client_id");
    const redirectUri = formData.get("redirect_uri");
    formData.forEach((value, key) => {
      formObject[key] = value;
    });

    try {
      const response = await fetch(
        `http://localhost:8080/oauth/login?client_id=${clientId}&redirect_uri=${redirectUri}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(formObject),
          credentials: "include",
        }
      );

      if (response.ok) {
        const data = await response.json();
        window.location.href =
          data.data.redirect_uri + "?code=" + data.data.auth_code;
      } else {
        const errorData = await response.json();
        responseElement.textContent = errorData.error || "Login Failed";
        responseElement.className = "error";
      }
    } catch (error) {
      responseElement.textContent = "Network Error: " + error.message;
      responseElement.className = "error";
    }
  });
});
