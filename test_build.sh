#!/bin/bash

# Test build script for OptiAssign

echo "Testing OptiAssign build..."

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first:"
    echo "   https://golang.org/dl/"
    exit 1
fi

# Install tools first
echo "Installing Go tools..."
go get -tool github.com/go-task/task/v3/cmd/task@latest

# Test Go build
echo "1. Testing Go build..."
if go tool task build; then
    echo "✅ Go build successful"
    go tool task clean
else
    echo "❌ Go build failed"
    exit 1
fi

# Test Go tests
echo "2. Running Go tests..."
if go tool task test; then
    echo "✅ Go tests passed"
else
    echo "❌ Go tests failed"
    exit 1
fi

# Test Docker build
echo "3. Testing Docker build..."
if go tool task docker-build; then
    echo "✅ Docker build successful"
    docker rmi optiassign:latest
else
    echo "❌ Docker build failed"
    exit 1
fi

echo "🎉 All tests passed! OptiAssign is ready for development."
echo ""
echo "🔧 Available Task commands:"
echo "   go tool task --list  # Show all available tasks"
echo "   go tool task dev     # Start development server"
echo "   go tool task test-coverage  # Run tests with coverage"
echo "   go tool task lint    # Run linter"