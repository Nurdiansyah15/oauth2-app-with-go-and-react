import React from "react";

function LoginPage({ onLogin }) {
  return (
    <div className="login-page">
      <h1 className="title">Welcome to App Two</h1>
      <p>Login to access your dashboard and manage your account.</p>
      <button className="login-button" onClick={onLogin}>
        Login
      </button>
    </div>
  );
}

export default LoginPage;
