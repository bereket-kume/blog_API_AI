# Blog API Test Suite

This document describes the comprehensive test suite for the Blog API components, following Clean Architecture principles.

## Test Structure

The test suite is organized into three main layers:

### 1. Unit Tests (with Mocks)
- **Location**: `usecases/blog_usecase_test.go`
- **Purpose**: Test business logic in isolation using mocked dependencies
- **Dependencies**: Mocked BlogRepository interface

### 2. Integration Tests (with Real Database)
- **Location**: `Infrastructure/repositories/blog_repository_test.go`
- **Purpose**: Test data persistence layer with real MongoDB
- **Dependencies**: Local MongoDB instance (Docker)

### 3. Controller Tests (with Mocks)
- **Location**: `Delivery/controllers/blog_controller_test.go`
- **Purpose**: Test HTTP endpoints and request/response handling
- **Dependencies**: Mocked BlogUseCase interface

## Test Components

### Mocks
- `mocks/blog_repository_mock.go` - Mock for BlogRepository interface
- `mocks/blog_usecase_mock.go` - Mock for BlogUseCase interface

### Test Utilities
- `test_helpers.go` - Common test data generators
- `test_config.go` - Test configuration and database utilities

## Prerequisites

### 1. Install Dependencies
```bash
# Install test dependencies
make test-deps

# Or manually
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/suite
```

### 2. MongoDB Setup
The integration tests require a local MongoDB instance. Make sure MongoDB is installed and running:

```bash
# Check if MongoDB is running
sudo systemctl status mongod

# Start MongoDB if not running
sudo systemctl start mongod

# Enable MongoDB to start on boot
sudo systemctl enable mongod

# Test MongoDB connection
mongo --eval "db.runCommand('ping')"
```

## Running Tests

### Using Makefile (Recommended)
```bash
# Run all tests
make test-all

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run specific component tests
make test-usecase
make test-controller
make test-repository

# Run with coverage
make coverage

# Check available commands
make help
```

### Using Go Commands
```bash
# Run all tests
go test ./...

# Run specific test files
go test ./usecases/blog_usecase_test.go
go test ./Delivery/controllers/blog_controller_test.go
go test ./Infrastructure/repositories/blog_repository_test.go

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

## Test Coverage

The test suite covers the following functionality:

### Blog UseCase Tests
- ✅ Create blog
- ✅ Get paginated blogs
- ✅ Get blog by ID
- ✅ Update blog
- ✅ Delete blog
- ✅ Search blogs
- ✅ Filter blogs
- ✅ Increment view count
- ✅ Update likes/dislikes
- ✅ Add comment
- ✅ Get comments

### Blog Repository Tests
- ✅ Database connection and setup
- ✅ CRUD operations
- ✅ Pagination
- ✅ Search functionality
- ✅ Filtering
- ✅ View count tracking
- ✅ Like/dislike tracking
- ✅ Comment management
- ✅ Error handling

### Blog Controller Tests
- ✅ HTTP request/response handling
- ✅ Authentication/authorization
- ✅ Input validation
- ✅ Error responses
- ✅ JSON serialization/deserialization
- ✅ Route handling

## Test Patterns

### 1. Table-Driven Tests
Tests use table-driven patterns for comprehensive coverage:

```go
tests := []struct {
    name        string
    input       interface{}
    setupMock   func(*mocks.BlogRepositoryMock)
    expectError bool
    expected    interface{}
}{
    // Test cases...
}
```

### 2. Setup and Teardown
Integration tests use proper setup and teardown:

```go
func (suite *BlogRepositoryTestSuite) SetupSuite() {
    // Connect to test database
}

func (suite *BlogRepositoryTestSuite) TearDownSuite() {
    // Clean up database
}

func (suite *BlogRepositoryTestSuite) SetupTest() {
    // Clear data before each test
}

func (suite *BlogRepositoryTestSuite) TearDownTest() {
    // Clear data after each test
}
```

### 3. Mock Expectations
Unit tests verify mock expectations:

```go
mockRepo.On("CreateBlog", mock.AnythingOfType("models.Blog")).Return(expectedBlog, nil)
// ... test execution ...
mockRepo.AssertExpectations(t)
```

## Test Data

### Test Helpers
The `test_helpers.go` file provides utilities for creating test data:

```go
helper := GetTestHelper()
blog := helper.CreateTestBlog()
blogs := helper.CreateTestBlogs(5)
comment := helper.CreateTestComment()
user := helper.CreateTestUser()
```

### Test Configuration
The `test_config.go` file manages test environment:

```go
config := GetTestConfig()
client, err := ConnectTestDB(config)
defer DisconnectTestDB(client, config)
```

## Environment Variables

You can customize test behavior with environment variables:

```bash
# MongoDB connection
export TEST_MONGO_URI="mongodb://localhost:27017"
export TEST_DB_NAME="blog_test_db"
export TEST_COLLECTION_NAME="blogs"
```

## Best Practices

### 1. Clean Architecture Compliance
- Unit tests use mocks to isolate business logic
- Integration tests test real database interactions
- Controller tests focus on HTTP concerns

### 2. Test Isolation
- Each test is independent
- Database is cleaned between tests
- Mocks are reset after each test

### 3. Comprehensive Coverage
- Happy path scenarios
- Error scenarios
- Edge cases
- Input validation

### 4. Readable Tests
- Descriptive test names
- Clear setup and assertions
- Meaningful test data

## Troubleshooting

### Common Issues

1. **MongoDB Connection Failed**
   ```bash
   # Check if MongoDB is running
   sudo systemctl status mongod
   
   # Start MongoDB if needed
   sudo systemctl start mongod
   
   # Check MongoDB logs
   sudo journalctl -u mongod -f
   ```

2. **Test Dependencies Missing**
   ```bash
   # Install dependencies
   make test-deps
   ```

3. **Port Already in Use**
   ```bash
   # Check what's using the port
   lsof -i :27017
   
   # Stop conflicting service
   sudo systemctl stop mongod
   ```

### Debug Mode
Run tests with verbose output for debugging:

```bash
make test-verbose
# or
go test -v ./...
```

## Continuous Integration

The test suite is designed to work in CI/CD environments:

```yaml
# Example GitHub Actions
- name: Run Tests
  run: |
    make test-deps
    make test-all-mongo
    make coverage
```

## Contributing

When adding new features:

1. Write tests first (TDD approach)
2. Follow existing test patterns
3. Ensure proper mock setup
4. Add integration tests for database operations
5. Update this documentation

## Test Reports

After running tests with coverage:

```bash
make coverage
```

Open `coverage/coverage.html` in your browser to view detailed coverage reports. 