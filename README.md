# RBAC User Management System (Go + React)

This is a modern, production-ready User Management System with Role-Based Access Control (RBAC).

## Tech Stack

- **Database**: MySQL with GORM.
- **Frontend**: React (Vite) with TypeScript, Tailwind CSS, and Zustand.
- **Environment**: Docker for containerized development.

## Features

- **Authentication**: JWT with Refresh Tokens.
- **Authorization**: Fine-grained RBAC middleware.
- **User Management**: CRUD, activation/deactivation, bulk actions.
- **Role Management**: Create, update, and assign permissions to roles.
- **Profile Management**: Self-service profile and password updates.
- **Dashboard**: Analytics on user activity and system health.

## Prerequisites

- Docker & Docker Compose
- Go (v1.21+)
- Node.js (v18+)
- pnpm (or npm/yarn)

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd rbac-system
```

### 2. Backend Setup

**a. Configure Environment Variables:**

Navigate to the `backend` directory and create a `.env` file from the example:

```bash
cd backend
cp .env.example .env
```

Update the `.env` file with your database credentials and JWT secret. These should match the values in `docker-compose.yml`.

```env
# Server
SERVER_PORT=8080
SERVER_ENV=development

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=rbac_system
DB_SSLMODE=false # Not applicable for MySQL, but keep for consistency if needed

# JWT
JWT_SECRET=your-super-secret-key
JWT_ACCESS_TOKEN_EXP_HOURS=1
JWT_REFRESH_TOKEN_EXP_HOURS=72

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173
```

**b. Start the Database:**

From the root directory, start the MySQL container:

```bash
docker-compose up -d
```

**c. Run the Backend Server:**

Navigate to the `backend` directory and run the application:

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

The backend server will be running at `http://localhost:8080`.

### 3. Frontend Setup

**a. Install Dependencies:**

Navigate to the `frontend` directory and install the required packages:

```bash
cd frontend
pnpm install
```

**b. Run the Frontend Development Server:**

```bash
pnpm run dev
```

The frontend application will be available at `http://localhost:5173`.

## Project Structure

```
rbac-system/
├── backend/      # Go Fiber application
├── frontend/     # React (Vite) application
├── docker-compose.yml # Docker setup for services
└── README.md     # This file
```

## API Documentation

API documentation can be generated using Swagger/OpenAPI (setup pending).

## Running Tests

To run backend tests:

```bash
cd backend
go test ./...
```

To run frontend tests:

```bash
cd frontend
pnpm run test
```
