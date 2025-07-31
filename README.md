# Blog API with AI

A modern Go-based blog API with AI integration using MongoDB Atlas, built with clean architecture principles.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- MongoDB (included in Docker setup)

### Setup
1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd blog_API_AI
   ```

2. **Configure environment:**
   ```bash
   cp env.example .env
   # Edit .env with your MongoDB Atlas credentials
   ```

3. **Run with Docker:**
   ```bash
   docker-compose up --build
   ```

4. **Verify installation:**
   ```bash
   curl http://localhost:8080/health
   ```

## 📚 Documentation

For detailed setup instructions, troubleshooting, and development guides, see:
- **[SETUP.md](SETUP.md)** - Comprehensive setup guide
- **[API Documentation](#api-endpoints)** - Available endpoints

## 🏗️ Architecture

This project follows Clean Architecture principles:

```
blog_API_AI/
├── Delivery/           # HTTP handlers and main application
├── Domain/            # Business logic and models
├── Infrastructure/    # Database, external services, utilities
├── Usecases/          # Application use cases
├── Dockerfile         # Docker configuration
├── docker-compose.yml # Docker Compose configuration
└── go.mod            # Go module file
```

## 🔌 Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/` | Welcome message |

**Note:** This is a configuration setup. API endpoints will be implemented in the respective controller files.

## 🛠️ Development

### Local Development
```bash
go mod download
go run Delivery/main.go
```

### Docker Commands
```bash
# Start application
docker-compose up -d

# View logs
docker-compose logs -f

# Stop application
docker-compose down

# Rebuild
docker-compose up --build
```

## 🔧 Configuration

Key environment variables:
- `MONGODB_URI` - MongoDB Atlas connection string
- `JWT_SECRET` - JWT signing secret
- `PORT` - Application port (default: 8080)
- `ENV` - Environment (development/production)

## 🚨 Troubleshooting

Common issues and solutions are documented in [SETUP.md](SETUP.md).

### Quick Fixes
- **Port conflict**: Change `PORT` in `.env`
- **MongoDB connection**: Verify Atlas credentials and IP whitelist
- **Docker issues**: Run `docker-compose down --rmi all` and rebuild

## 🔒 Security

- Environment variables for sensitive data
- Non-root Docker container
- CORS configuration
- JWT authentication ready

## 📄 License

[Add your license information here]

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

---

**Need help?** Check the [SETUP.md](SETUP.md) for detailed instructions and troubleshooting.
