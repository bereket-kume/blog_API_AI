# Blog API with AI - Setup Guide

This guide will help you set up and run the Blog API with AI application on your machine using local MongoDB.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Docker** (version 20.10 or higher)
- **Docker Compose** (version 2.0 or higher)
- **Git** (for cloning the repository)

### Installing Docker

#### Ubuntu/Debian:
```bash
# Update package index
sudo apt-get update

# Install prerequisites
sudo apt-get install apt-transport-https ca-certificates curl gnupg lsb-release

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up stable repository
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Add user to docker group (optional, to run docker without sudo)
sudo usermod -aG docker $USER
```

#### macOS:
```bash
# Install Docker Desktop
# Download from: https://www.docker.com/products/docker-desktop
# Follow the installation wizard
```

#### Windows:
```bash
# Install Docker Desktop
# Download from: https://www.docker.com/products/docker-desktop
# Follow the installation wizard
```

## Application Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd blog_API_AI
```

### 2. Configure Environment Variables
```bash
# Copy the environment template
cp env.example .env

# The .env file is already configured for local MongoDB
# No additional configuration needed
```

The `.env` file contains:
```env
# MongoDB Configuration (Local)
MONGODB_URI=mongodb://admin:password123@mongodb:27017/blog_db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Application Configuration
PORT=8080
ENV=production

# Local MongoDB Configuration
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=password123
MONGO_DATABASE=blog_db
```

### 3. Build and Run with Docker
```bash
# Build and start the application
docker-compose up --build

# Or run in detached mode (background)
docker-compose up -d --build
```

This will start:
- **blog-api**: Your Go application
- **mongodb**: Local MongoDB database

### 4. Verify Installation
```bash
# Check if containers are running
docker-compose ps

# Test the health endpoint
curl http://localhost:8080/health

# Test the root endpoint
curl http://localhost:8080/
```

## Development Setup (Optional)

If you want to run the application locally without Docker:

### 1. Install Go
```bash
# Ubuntu/Debian
sudo apt-get install golang-go

# macOS
brew install go

# Windows
# Download from: https://golang.org/dl/
```

### 2. Install MongoDB Locally
```bash
# Ubuntu/Debian
sudo apt-get install mongodb

# macOS
brew install mongodb/brew/mongodb-community

# Windows
# Download from: https://www.mongodb.com/try/download/community
```

### 3. Install Dependencies
```bash
go mod download
```

### 4. Run Locally
```bash
# Start MongoDB (if not running as a service)
mongod

# In another terminal, run the application
go run Delivery/main.go
```

## API Endpoints

Once the application is running, you can access the following endpoints:

### Health Check
- `GET /health` - Check application health

### Root
- `GET /` - Welcome message and API information

**Note:** This is a configuration setup. The actual API endpoints (blogs, users, authentication) will be implemented in the respective controller files following the clean architecture pattern.

## Docker Commands

### Basic Commands
```bash
# Start the application
docker-compose up

# Start in background
docker-compose up -d

# Stop the application
docker-compose down

# View logs
docker-compose logs

# View logs in real-time
docker-compose logs -f

# Rebuild and start
docker-compose up --build

# Stop and remove containers, networks, and images
docker-compose down --rmi all --volumes --remove-orphans
```

### Container Management
```bash
# List running containers
docker-compose ps

# Execute commands in the container
docker-compose exec blog-api sh

# View container logs
docker-compose logs blog-api

# Restart the service
docker-compose restart blog-api

# Access MongoDB shell
docker-compose exec mongodb mongosh -u admin -p password123
```

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Check what's using port 8080
sudo lsof -i :8080

# Kill the process or change the port in .env
PORT=8081
```

#### 2. MongoDB Connection Failed
- Check if MongoDB container is running: `docker-compose ps`
- View MongoDB logs: `docker-compose logs mongodb`
- Ensure the connection string in `.env` is correct
- Try restarting the containers: `docker-compose restart`

#### 3. Docker Build Fails
```bash
# Clean Docker cache
docker system prune -a

# Rebuild without cache
docker-compose build --no-cache
```

#### 4. Permission Issues
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and log back in, or run:
newgrp docker
```

### Health Check
The application includes a health check endpoint at `/health` that returns:
```json
{
  "status": "healthy",
  "timestamp": "2025-07-30T14:19:11.298407984Z",
  "service": "blog-api"
}
```

### Logs
```bash
# View application logs
docker-compose logs -f blog-api

# View MongoDB logs
docker-compose logs -f mongodb

# View all logs
docker-compose logs -f
```

## MongoDB Management

### Access MongoDB Shell
```bash
# Connect to MongoDB container
docker-compose exec mongodb mongosh -u admin -p password123

# List databases
show dbs

# Use blog database
use blog_db

# List collections
show collections
```

### Backup and Restore
```bash
# Backup database
docker-compose exec mongodb mongodump -u admin -p password123 --db blog_db --out /backup

# Restore database
docker-compose exec mongodb mongorestore -u admin -p password123 --db blog_db /backup/blog_db
```

## Security Notes

- Never commit your `.env` file to version control
- Use strong JWT secrets in production
- Consider using secrets management for production deployments
- Regularly update dependencies
- Change default MongoDB credentials in production

## Production Deployment

For production deployment:

1. **Environment Variables**: Use proper secrets management
2. **SSL/TLS**: Configure HTTPS
3. **Load Balancer**: Set up proper load balancing
4. **Monitoring**: Add application monitoring
5. **Backup**: Configure MongoDB backups
6. **Security**: Implement proper authentication and authorization
7. **Database**: Consider using managed MongoDB service

## Support

If you encounter any issues:

1. Check the troubleshooting section above
2. Review the application logs
3. Verify your Docker installation
4. Ensure all prerequisites are installed correctly

## Project Structure

```
blog_API_AI/
├── Delivery/           # HTTP handlers and main application
│   └── main.go        # Application entry point
├── Domain/            # Business logic and models
│   ├── interfaces/    # Repository interfaces
│   └── models/        # Data models
├── Infrastructure/    # Database, external services, utilities
│   ├── database/      # Database connections
│   ├── repositories/  # Data access layer
│   ├── services/      # External services (AI, JWT)
│   └── utils/         # Utility functions
├── Usecases/          # Application use cases
├── Dockerfile         # Docker configuration
├── docker-compose.yml # Docker Compose configuration
├── go.mod            # Go module file
├── .env              # Environment variables (create from env.example)
└── README.md         # Project documentation
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Add your license information here]
