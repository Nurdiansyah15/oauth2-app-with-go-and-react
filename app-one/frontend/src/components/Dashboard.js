import React, { useState } from "react";

function Dashboard({ user, onLogout }) {
  const [showDropdown, setShowDropdown] = useState(false);

  const toggleDropdown = () => {
    setShowDropdown(!showDropdown);
  };

  const navigateToApp = (url) => {
    window.location.href = url;
  };

  return (
    <div className="dashboard">
      {/* Navbar */}
      <nav className="navbar">
        <h2 className="navbar-title">App One</h2>
        <div className="dropdown">
          <button className="dropdown-button" onClick={toggleDropdown}>
            Switch App
          </button>
          {showDropdown && (
            <ul className="dropdown-menu">
              <li onClick={() => navigateToApp("http://localhost:3001")}>
                App Two
              </li>
              <li onClick={() => navigateToApp("http://localhost:3002")}>
                App Three
              </li>
            </ul>
          )}
        </div>
      </nav>

      {/* Dashboard Content */}
      <div className="dashboard-content">
        <h1 className="dashboard-title">Hello, {user?.username || "User"}!</h1>
        <p>Welcome to your dashboard. Manage your activities here.</p>
        <button className="logout-button" onClick={onLogout}>
          Logout
        </button>
      </div>
    </div>
  );
}

export default Dashboard;
