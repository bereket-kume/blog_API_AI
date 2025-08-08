#!/bin/bash

# Blog API Test Runner Script
# This script provides an easy way to run different types of tests

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if MongoDB is running
check_mongo() {
	if pgrep -x "mongod" > /dev/null; then
		print_success "Local MongoDB is running"
		return 0
	else
		print_warning "Local MongoDB is not running"
		return 1
	fi
}

# Function to check MongoDB connection
check_mongo_connection() {
	print_status "Checking MongoDB connection..."
	# Try mongosh first (newer MongoDB versions)
	if command -v mongosh > /dev/null 2>&1; then
		if mongosh --eval "db.runCommand('ping')" > /dev/null 2>&1; then
			print_success "MongoDB connection successful (mongosh)"
			return 0
		fi
	fi
	# Try mongo (older MongoDB versions)
	if command -v mongo > /dev/null 2>&1; then
		if mongo --eval "db.runCommand('ping')" > /dev/null 2>&1; then
			print_success "MongoDB connection successful (mongo)"
			return 0
		fi
	fi
	# Try direct connection test
	if nc -z localhost 27017 2>/dev/null; then
		print_success "MongoDB port is accessible"
		return 0
	fi
	print_error "MongoDB connection failed"
	return 1
}

# Function to install test dependencies
install_deps() {
    print_status "Installing test dependencies..."
    go mod tidy
    print_success "Dependencies installed"
}

# Function to run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    if go test -v -timeout 30s ./usecases/... ./Delivery/controllers/...; then
        print_success "Unit tests passed"
    else
        print_error "Unit tests failed"
        exit 1
    fi
}

# Function to run integration tests
run_integration_tests() {
    print_status "Running integration tests..."
    if go test -v -timeout 30s ./Infrastructure/repositories/...; then
        print_success "Integration tests passed"
    else
        print_error "Integration tests failed"
        exit 1
    fi
}

# Function to run all tests
run_all_tests() {
    print_status "Running all tests..."
    if go test -v -timeout 30s ./...; then
        print_success "All tests passed"
    else
        print_error "Some tests failed"
        exit 1
    fi
}

# Function to run tests with coverage
run_coverage_tests() {
    print_status "Running tests with coverage..."
    mkdir -p coverage
    if go test -v -coverprofile=coverage/coverage.out -timeout 30s ./...; then
        go tool cover -html=coverage/coverage.out -o coverage/coverage.html
        print_success "Coverage report generated at coverage/coverage.html"
    else
        print_error "Tests failed during coverage run"
        exit 1
    fi
}

# Function to run specific test file
run_specific_test() {
    local test_file=$1
    print_status "Running test: $test_file"
    if go test -v -timeout 30s "$test_file"; then
        print_success "Test passed: $test_file"
    else
        print_error "Test failed: $test_file"
        exit 1
    fi
}

# Function to clean up
cleanup() {
    print_status "Cleaning up..."
    go clean -testcache
    print_success "Cleanup completed"
}

# Function to show help
show_help() {
    echo "Blog API Test Runner"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  unit              Run unit tests only"
    echo "  integration       Run integration tests only"
    echo "  all               Run all tests"
    echo "  coverage          Run tests with coverage report"
    echo "  usecase           Run blog usecase tests"
    echo "  controller        Run blog controller tests"
    echo "  repository        Run blog repository tests"
    echo "  deps              Install test dependencies"
    echo "  mongo             Check local MongoDB status"
    echo "  clean             Clean up test artifacts"
    echo "  help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 unit           # Run unit tests"
    echo "  $0 integration    # Run integration tests"
    echo "  $0 all            # Run all tests"
    echo "  $0 coverage       # Run tests with coverage"
}

# Main script logic
case "${1:-help}" in
    "unit")
        install_deps
        run_unit_tests
        ;;
    "integration")
        install_deps
        if ! check_mongo; then
            print_error "Local MongoDB is not running. Please start MongoDB manually."
            print_status "You can start MongoDB with: sudo systemctl start mongod"
            exit 1
        fi
        if ! check_mongo_connection; then
            print_error "Cannot connect to MongoDB. Please check your MongoDB configuration."
            exit 1
        fi
        run_integration_tests
        ;;
    "all")
        install_deps
        if ! check_mongo; then
            print_error "Local MongoDB is not running. Please start MongoDB manually."
            print_status "You can start MongoDB with: sudo systemctl start mongod"
            exit 1
        fi
        if ! check_mongo_connection; then
            print_error "Cannot connect to MongoDB. Please check your MongoDB configuration."
            exit 1
        fi
        run_all_tests
        ;;
    "coverage")
        install_deps
        if ! check_mongo; then
            print_error "Local MongoDB is not running. Please start MongoDB manually."
            print_status "You can start MongoDB with: sudo systemctl start mongod"
            exit 1
        fi
        if ! check_mongo_connection; then
            print_error "Cannot connect to MongoDB. Please check your MongoDB configuration."
            exit 1
        fi
        run_coverage_tests
        ;;
    "usecase")
        install_deps
        run_specific_test "./usecases/blog_usecase_test.go"
        ;;
    "controller")
        install_deps
        run_specific_test "./Delivery/controllers/blog_controller_test.go"
        ;;
    "repository")
        install_deps
        if ! check_mongo; then
            print_error "Local MongoDB is not running. Please start MongoDB manually."
            print_status "You can start MongoDB with: sudo systemctl start mongod"
            exit 1
        fi
        if ! check_mongo_connection; then
            print_error "Cannot connect to MongoDB. Please check your MongoDB configuration."
            exit 1
        fi
        run_specific_test "./Infrastructure/repositories/blog_repository_test.go"
        ;;
    "deps")
        install_deps
        ;;
    "mongo")
        if ! check_mongo; then
            print_error "Local MongoDB is not running."
            print_status "To start MongoDB, run: sudo systemctl start mongod"
            print_status "To enable MongoDB on boot, run: sudo systemctl enable mongod"
            exit 1
        fi
        if ! check_mongo_connection; then
            print_error "Cannot connect to MongoDB."
            print_status "Please check your MongoDB configuration and ensure it's running on localhost:27017"
            exit 1
        fi
        print_success "MongoDB is running and accessible"
        ;;
    "clean")
        cleanup
        ;;
    "help"|*)
        show_help
        ;;
esac 