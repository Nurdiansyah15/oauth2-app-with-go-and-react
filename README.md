# OAuth Application with Multiple Frontends and Backends

This repository hosts a multi-container OAuth application with two frontend apps and two backend apps, along with an authentication server. It leverages Docker Compose to manage the containers for smooth communication and deployment.

## Table of Contents
- [Overview](#overview)
- [Technologies Used](#technologies-used)
- [Folder Structure](#folder-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Available Endpoints](#available-endpoints)
- [Author](#author)

---

## Overview

The application is designed to:
- Authenticate users via an OAuth server.
- Allow communication between two separate frontend-backend pairs.
- Use PostgreSQL for data persistence.

Each service is containerized using Docker, and communication is managed within a shared Docker network.

### Services

1. **PostgreSQL**: Database for storing user and token data.
2. **Auth Server**: Manages OAuth authentication and token issuance.
3. **Frontend App One**: A user-facing application (App One).
4. **Backend App One**: API server for Frontend App One.
5. **Frontend App Two**: A user-facing application (App Two).
6. **Backend App Two**: API server for Frontend App Two.

---

## Technologies Used

- Docker & Docker Compose
- PostgreSQL 16
- OAuth 2.0
- Gin (for backend apps)
- React (for frontend apps)

---

## Folder Structure

```plaintext
.
├── auth-server/            # Auth server source code
├── app-one/
│   ├── frontend/           # Frontend App One code
│   └── backend/            # Backend App One code
├── app-two/
│   ├── frontend/           # Frontend App Two code
│   └── backend/            # Backend App Two code
├── docker-compose.yml      # Docker Compose configuration
└── README.md               # Project documentation
```

---

## Getting Started

### Prerequisites

- Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).

### Steps to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/Nurdiansyah15/oauth2-app-with-go-and-react.git
   cd oauth2-app-with-go-and-react
   ```

2. Build and start all containers:
   ```bash
   docker-compose up --build
   ```

3. Access the applications:
   - **Frontend App One**: [http://localhost:3000](http://localhost:3000)
   - **Frontend App Two**: [http://localhost:3001](http://localhost:3001)
   - **Auth Server**: [http://localhost:8080](http://localhost:8080)

4. Login
   You can using credential in seeder auth-server.
   Example : Nurdiansyah (nurdiansyah)


---

## Environment Variables

Each service requires specific environment variables. Below are the key ones used in the application:

### Auth Server
- `DB_HOST`: Hostname of the PostgreSQL container (default: `postgres-container`).
- `DB_PORT`: Port of the PostgreSQL database (default: `5432`).
- `DB_USER`: PostgreSQL username (default: `postgres`).
- `DB_PASS`: PostgreSQL password (default: `postgres`).
- `DB_NAME`: Database name (default: `oauth_db`).

### Frontend and Backend Apps
No environment variables required for the frontend and backend apps in this setup.

---

## Available Endpoints

### Auth Server
- `GET /login`: Displays the login page.
- `POST /oauth/login`: Handles user login and authentication.
- `GET /oauth/authorize`: Handles the OAuth authorization process.
- `POST /oauth/token`: Exchanges authorization codes for access tokens.
- `GET /oauth/validate`: Validates access tokens.
- `GET /oauth/logout`: Logs out the authenticated user and invalidates the session.

### Backend App One
- `POST /auth/callback`: Handles OAuth callback and token exchange.
- `GET /api/me`: Proxies request to the auth server.

### Backend App Two
- Extend as needed for App Two's functionality.

---

## Troubleshooting

1. **Database connection issues:**
   Ensure the `postgres` container is running and properly exposed to other services.

2. **Container communication issues:**
   Verify services are on the same Docker network (`app-network`).

3. **Connection refused errors:**
   Use service names (e.g., `auth-server`) instead of `localhost` for inter-container communication.

4. **Debugging within containers:**
   ```bash
   docker exec -it <container_name> sh
   ```

---

## Author

**Nurdiansyah**  
- GitHub: [Nurdiansyah15](https://github.com/Nurdiansyah15)  

Feel free to contribute or raise issues in this repository!
