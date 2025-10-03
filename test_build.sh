#!/bin/bash

# Test build script for OptiAssign

echo "Testing OptiAssign build..."

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first:"
    echo "   https://golang.org/dl/"
    exit 1
fi

# Test Go build
echo "1. Testing Go build..."
if go run github.com/go-task/task/v3/cmd/task build; then
    echo "✅ Go build successful"
    go run github.com/go-task/task/v3/cmd/task clean
else
    echo "❌ Go build failed"
    exit 1
fi

# Test Go tests
echo "2. Running Go tests..."
if go run github.com/go-task/task/v3/cmd/task test; then
    echo "✅ Go tests passed"
else
    echo "❌ Go tests failed"
    exit 1
fi

# Test Docker build
echo "3. Testing Docker build..."
if go run github.com/go-task/task/v3/cmd/task docker-build; then
    echo "✅ Docker build successful"
    docker rmi optiassign:latest
else
    echo "❌ Docker build failed"
    exit 1
fi

echo "🎉 All tests passed! OptiAssign is ready for development."
echo ""
echo "🔧 Available Task commands:"
echo "   go run github.com/go-task/task/v3/cmd/task --list  # Show all available tasks"
echo "   go run github.com/go-task/task/v3/cmd/task dev     # Start development server"
echo "   go run github.com/go-task/task/v3/cmd/task test-coverage  # Run tests with coverage"
echo "   go run github.com/go-task/task/v3/cmd/task lint    # Run linter"