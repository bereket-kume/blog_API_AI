# Blog API

A RESTful API built with **Go**, **Gin**, and **MongoDB** that handles user authentication, authorization, and role-based access control with JWT tokens (access & refresh).

---

## 🚀 Features

- **User Registration** (automatic role assignment: first user becomes admin)
- **User Login** with access & refresh tokens
- **Refresh Token** endpoint to get a new access token
- **Role-based Authorization** (admin & superadmin)
- **Password hashing** with bcrypt
- **MongoDB Atlas** (or local MongoDB) integration
- **Clean Architecture** with interfaces, usecases, and repositories
- **Unit & integration tests** with `testify`

---

## 🛠️ Tech Stack

- **Language**: Go 1.20+
- **Web Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT (access & refresh tokens)
- **Hashing**: bcrypt
- **Testing**: Testify
- **Dependency Injection**: Manual

---

## 📦 Installation

```bash
# 1. Clone the repository
git clone https://github.com/your-username/blog-api.git
cd blog-api

# 2. Install dependencies
go mod tidy

# 3. Set environment variables
export MONGODB_URI="mongodb://localhost:27017" # or MongoDB Atlas URI
export ACCESS_SECRET="your_access_secret"
export REFRESH_SECRET="your_refresh_secret"

# 4. Run the server
go run Delivery/main.go
```

---

## 🔑 API Endpoints

### Auth

#### 1. Register

**POST** `/register`

**Request:**

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:**

```json
{
  "message": "User registered successfully"
}
```

---

#### 2. Login

**POST** `/login`

**Request:**

```json
{
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:**

```json
{
  "access_token": "jwt_access_token_here",
  "refresh_token": "jwt_refresh_token_here"
}
```

---

#### 3. Refresh Token

**POST** `/refresh`

**Request:**

```json
{
  "refresh_token": "jwt_refresh_token_here"
}
```

**Response:**

```json
{
  "access_token": "new_jwt_access_token_here"
}
```

---

### Admin

#### Promote User

**POST** `/api/promote`

**Headers:**

```
Authorization: Bearer <admin_or_superadmin_access_token>
```

**Request:**

```json
{
  "email": "user@example.com"
}
```

**Response:**

```json
{
  "message": "User promoted to admin"
}
```

---

### Superadmin

#### Demote Admin

**POST** `/api/demote`

**Headers:**

```
Authorization: Bearer <superadmin_access_token>
```

**Request:**

```json
{
  "email": "admin@example.com"
}
```

**Response:**

```json
{
  "message": "Admin demoted to user"
}
```

---

## 🧪 Running Tests

```bash
# Run all tests
go test ./...

# Run specific test file
go test ./repositories -v
```

---

## 📂 Project Structure

```
blog-api/
│
├── Delivery/             # Entry points & HTTP layer
│   ├── controllers/      # Gin handlers
│   ├── middlewares/      # Auth middleware
│   ├── routers/          # Routes setup
│   └── main.go           # Application entry
│
├── Domain/               # Entities & interfaces
│   ├── interfaces/       # Interface contracts
│   └── models/           # Domain models
│
├── Infrastructure/       # Database & external services
│   ├── database/         # MongoDB connection
│   ├── db_models/        # DB-specific models
│   └── repositories/     # Repository implementations
│
├── services/             # JWT & hashing services
├── usecases/             # Business logic
├── mocks/                # Test mocks
└── go.mod
```

---

## 🛡️ Security Notes

- Passwords are hashed with bcrypt (cost 14)
- Refresh tokens are stored hashed in the database
- Expired refresh tokens are deleted
- Role checks are enforced in middleware

---
