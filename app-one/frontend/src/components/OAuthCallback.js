import React, { useEffect } from "react";

function OAuthCallback() {
  useEffect(() => {
    const handleCallback = async () => {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");

      if (code) {
        try {
          const response = await fetch("http://localhost:8081/auth/callback", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ code }),
          });

          if (response.ok) {
            const data = await response.json();
            localStorage.setItem("token", data.data.access_token);
            window.location.href = "/";
          }
        } catch (error) {
          console.error("Auth callback failed:", error);
          window.location.href = "/login";
        }
      }
    };

    handleCallback();
  }, []);

  return <div className="callback">Processing your login...</div>;
}

export default OAuthCallback;
