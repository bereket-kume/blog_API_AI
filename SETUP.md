# Blog API with AI - Setup Guide

This guide will help you set up and run the Blog API with AI application on your machine.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Docker** (version 20.10 or higher)
- **Docker Compose** (version 2.0 or higher)
- **Git** (for cloning the repository)
- **MongoDB Atlas account** (free tier available)

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

## MongoDB Atlas Setup

### 1. Create MongoDB Atlas Account
1. Go to [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Click "Try Free" and create an account
3. Choose the free tier (M0)

### 2. Create a Cluster
1. Click "Build a Database"
2. Choose "FREE" tier
3. Select your preferred cloud provider and region
4. Click "Create"

### 3. Configure Database Access
1. Go to "Database Access" in the left sidebar
2. Click "Add New Database User"
3. Choose "Password" authentication
4. Create a username and password (save these!)
5. Set privileges to "Read and write to any database"
6. Click "Add User"

### 4. Configure Network Access
1. Go to "Network Access" in the left sidebar
2. Click "Add IP Address"
3. For development, click "Allow Access from Anywhere" (0.0.0.0/0)
4. Click "Confirm"

### 5. Get Connection String
1. Go to "Database" in the left sidebar
2. Click "Connect" on your cluster
3. Choose "Connect your application"
4. Copy the connection string

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

# Edit the .env file with your MongoDB Atlas credentials
nano .env
```

Update the `.env` file with your MongoDB Atlas connection string:
```env
# MongoDB Atlas Configuration
MONGODB_URI=mongodb+srv://your_username:your_password@your_cluster.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Application Configuration
PORT=8080
ENV=production
```

### 3. Build and Run with Docker
```bash
# Build and start the application
docker-compose up --build

# Or run in detached mode (background)
docker-compose up -d --build
```

### 4. Verify Installation
```bash
# Check if the container is running
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

### 2. Install Dependencies
```bash
go mod download
```

### 3. Run Locally
```bash
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
- Verify your MongoDB Atlas connection string
- Check if your IP is whitelisted in Atlas
- Ensure database user has correct permissions
- Verify username and password are correct

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

# View specific log levels
docker-compose logs -f --tail=100 blog-api
```

## Security Notes

- Never commit your `.env` file to version control
- Use strong JWT secrets in production
- Configure proper network access in MongoDB Atlas
- Consider using secrets management for production deployments
- Regularly update dependencies

## Production Deployment

For production deployment:

1. **Environment Variables**: Use proper secrets management
2. **SSL/TLS**: Configure HTTPS
3. **Load Balancer**: Set up proper load balancing
4. **Monitoring**: Add application monitoring
5. **Backup**: Configure MongoDB Atlas backups
6. **Security**: Implement proper authentication and authorization

## Support

If you encounter any issues:

1. Check the troubleshooting section above
2. Review the application logs
3. Verify your MongoDB Atlas configuration
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