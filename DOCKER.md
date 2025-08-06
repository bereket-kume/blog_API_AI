# Docker Setup

This project includes a simple Docker configuration to run both the Go API and MongoDB database.

## Quick Start

1. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

2. **Run in background:**
   ```bash
   docker-compose up -d --build
   ```

3. **Stop the services:**
   ```bash
   docker-compose down
   ```

4. **View logs:**
   ```bash
   docker-compose logs -f app
   ```

## Services

- **App**: Go API running on port 8080
- **MongoDB**: Database running on port 27017

## Environment Variables

The application uses the following environment variables (configured in docker-compose.yml):

- `MONGODB_URI`: MongoDB connection string
- `JWT_SECRET`: Secret key for JWT tokens
- `PORT`: Application port (8080)
- `ENV`: Environment (production)

## Data Persistence

MongoDB data is persisted in a Docker volume named `mongodb_data`.

## Development

To rebuild the application after code changes:

```bash
docker-compose up --build
```

## Clean Up

To remove all containers, networks, and volumes:

```bash
docker-compose down -v
``` 