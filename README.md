# Blog API with AI Integration

A modern, scalable blog API built with Go following Clean Architecture principles, featuring AI-powered content suggestions, intelligent recommendations, and comprehensive user management.

## 🚀 Features

### Core Blog Functionality
- **Blog Management**: Create, read, update, and delete blog posts
- **Comment System**: Full CRUD operations for blog comments
- **Like/Dislike System**: User engagement tracking
- **Search & Filtering**: Advanced content discovery with pagination
- **Tag System**: Categorized content organization

### AI-Powered Features
- **Content Suggestions**: AI-generated blog ideas and content recommendations
- **Smart Recommendations**: Personalized content based on user behavior
- **Content Discovery**: Trending, popular, and new content algorithms
- **Similar Content**: Find related posts using AI similarity analysis

### User Management
- **Authentication**: JWT-based secure authentication with refresh tokens
- **Role-Based Access Control**: User, Admin, and Superadmin roles
- **Email Verification**: Secure email verification system
- **Password Reset**: Secure password recovery via email
- **User Profiles**: Rich user profiles with bio and contact information

### Recommendation Engine
- **Behavioral Tracking**: Monitor user interactions (views, likes, comments)
- **Interest Analysis**: Build user interest profiles
- **Personalized Feed**: AI-driven content recommendations
- **Performance Analytics**: Track recommendation effectiveness

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
blog_API_AI/
├── Domain/           # Business logic and interfaces
│   ├── interfaces/   # Abstract interfaces
│   └── models/       # Domain entities
├── usecases/         # Business use cases
├── Infrastructure/   # External concerns
│   ├── database/     # MongoDB connection
│   ├── repositories/ # Data access layer
│   ├── services/     # External services
│   └── utils/        # Utility functions
├── Delivery/         # HTTP layer
│   ├── controllers/  # Request handlers
│   ├── middlewares/  # HTTP middleware
│   └── routers/      # Route definitions
└── mocks/            # Test mocks
```

### Architecture Benefits
- **Testability**: Easy to mock dependencies for unit testing
- **Maintainability**: Clear separation of concerns
- **Scalability**: Modular design allows easy scaling
- **Flexibility**: Easy to swap implementations

## 🛠️ Technology Stack

- **Language**: Go 1.23+
- **Framework**: Gin (HTTP web framework)
- **Database**: MongoDB with official Go driver
- **Authentication**: JWT with secure token management
- **Email**: Brevo SMTP integration
- **Testing**: Testify framework with comprehensive mocking
- **AI Integration**: External AI service for content suggestions

## 📋 Prerequisites

- Go 1.23 or higher
- MongoDB (local or Atlas)
- Git

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd blog_API_AI
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Environment Configuration
```bash
cp env.example .env
# Edit .env with your configuration
```

### 4. Database Setup
#### Option A: Local MongoDB
```bash
# Start MongoDB locally
sudo systemctl start mongod

# Or using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

#### Option B: MongoDB Atlas
- Create a MongoDB Atlas cluster
- Update `MONGODB_URI` in your `.env` file

### 5. Run the Application
```bash
go run Delivery/main.go
```

The API will be available at `http://localhost:8080`

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017/blog_db` |
| `JWT_SECRET` | JWT signing secret | Required |
| `PORT` | Server port | `8080` |
| `ENV` | Environment (development/production) | `development` |
| `BREVO_SMTP_HOST` | SMTP server host | `smtp-relay.brevo.com` |
| `BREVO_SMTP_PORT` | SMTP server port | `587` |
| `BREVO_SMTP_USERNAME` | SMTP username | Required |
| `BREVO_SMTP_PASSWORD` | SMTP password | Required |
| `FROM_EMAIL` | Sender email address | Required |
| `FRONTEND_URL` | Frontend application URL | `http://localhost:3000` |

## 📚 API Documentation


### Key Endpoints

#### Authentication
- `POST /register` - User registration
- `POST /login` - User login
- `POST /refresh` - Refresh access token
- `POST /logout` - User logout
- `GET /verify-email` - Email verification
- `POST /forgot-password` - Password reset request
- `GET /reset-password` - Password reset

#### Blogs (Public)
- `GET /blogs` - Get paginated blogs
- `GET /blogs/search` - Search blogs
- `GET /blogs/filter` - Filter blogs
- `GET /blogs/:id` - Get blog by ID
- `GET /blogs/:id/comments` - Get blog comments

#### Blogs (Authenticated)
- `POST /api/blogs` - Create blog
- `PUT /api/blogs/:id` - Update blog
- `DELETE /api/blogs/:id` - Delete blog
- `POST /api/blogs/:id/comments` - Add comment
- `POST /api/blogs/:id/like` - Like blog
- `POST /api/blogs/:id/unlike` - Unlike blog
- `POST /api/blogs/:id/dislike` - Dislike blog

#### AI Features (Authenticated)
- `POST /api/ai/suggestions` - Generate AI suggestions
- `POST /api/ai/ideas` - Generate content ideas
- `POST /api/ai/save` - Save AI suggestion
- `GET /api/ai/suggestions` - Get AI suggestions
- `POST /api/ai/suggestions/:id/convert-to-draft` - Convert to draft

#### Recommendations
- `GET /recommendations/trending` - Get trending content
- `GET /recommendations/popular` - Get popular content
- `GET /api/recommendations/personal` - Get personalized recommendations
- `POST /api/recommendations/track` - Track user behavior

#### User Management
- `GET /api/user/profile` - Get user profile
- `PUT /api/user/profile` - Update user profile
- `POST /api/admin/promote` - Promote user (Admin only)
- `POST /api/superadmin/demote` - Demote user (Superadmin only)

## 🧪 Testing

### Run All Tests
```bash
make test-all
```

### Run Specific Test Types
```bash
make test-unit          # Unit tests only
make test-integration   # Integration tests only
make test-usecase       # Blog usecase tests
make test-controller    # Controller tests
make test-repository    # Repository tests
```

### Test Coverage
```bash
make coverage
# Opens coverage.html in browser
```

### Test with MongoDB Check
```bash
make test-all-mongo     # Checks MongoDB before running tests
```

### Available Test Commands
```bash
make help               # Show all available test commands
make test-verbose       # Verbose test output
make test-race          # Race condition detection
make test-short         # Short timeout for CI/CD
make clean              # Clean test artifacts
```

## 🔧 Development

### Project Structure
```
├── Domain/                 # Business domain layer
│   ├── interfaces/         # Abstract interfaces
│   └── models/            # Domain entities
├── usecases/              # Business logic implementation
├── Infrastructure/        # External concerns
│   ├── database/          # Database connection
│   ├── repositories/      # Data access implementation
│   ├── services/          # External service implementations
│   └── utils/             # Utility functions
├── Delivery/              # HTTP delivery layer
│   ├── controllers/       # HTTP request handlers
│   ├── middlewares/       # HTTP middleware
│   └── routers/           # Route definitions
├── mocks/                 # Test mocks
└── tests/                 # Test files
```

### Adding New Features

1. **Define Domain Models** in `Domain/models/`
2. **Create Interfaces** in `Domain/interfaces/`
3. **Implement Use Cases** in `usecases/`
4. **Add Repositories** in `Infrastructure/repositories/`
5. **Create Controllers** in `Delivery/controllers/`
6. **Update Routes** in `Delivery/routers/`
7. **Write Tests** for all layers

### Code Style
- Follow Go standard formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Follow Clean Architecture principles

## 🚀 Deployment

### Production Considerations
- Set `ENV=production` in environment
- Use strong `JWT_SECRET`
- Configure MongoDB Atlas for production
- Set up proper CORS policies
- Use HTTPS in production
- Configure proper logging

### Docker Deployment
```bash
# Build Docker image
docker build -t blog-api .

# Run container
docker run -p 8080:8080 --env-file .env blog-api
```

### Environment-Specific Configs
- **Development**: Local MongoDB, debug logging
- **Staging**: MongoDB Atlas, moderate logging
- **Production**: MongoDB Atlas, minimal logging, HTTPS

## 📊 Performance & Monitoring

### Database Optimization
- MongoDB indexes on frequently queried fields
- Connection pooling
- Query optimization

### API Performance
- Pagination for large datasets
- Caching strategies
- Rate limiting (can be added)

### Monitoring
- Health check endpoint (`/health`)
- Structured logging
- Error tracking

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: Bcrypt password hashing
- **Role-Based Access Control**: Granular permission system
- **Input Validation**: Comprehensive input sanitization
- **CORS Configuration**: Configurable cross-origin policies
- **Secure Headers**: Security-focused HTTP headers

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Clean Architecture principles
- Write comprehensive tests
- Update documentation
- Follow Go coding standards

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

- **Issues**: Create GitHub issues for bugs or feature requests
- **Documentation**: Check the [API docs](https://documenter.getpostman.com/view/38774125/2sB3BEnVMb#6a2d9476-e662-4bbb-9ec7-882c0af10d4c)
- **Testing**: Use the comprehensive test suite


---

**Built with ❤️ using Go and Clean Architecture principles**
