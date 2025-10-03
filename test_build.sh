#!/bin/bash

# Test build script for OptiAssign

echo "Testing OptiAssign build..."

# Check if Task is installed
if ! command -v task &> /dev/null; then
    echo "❌ Task is not installed. Please install it first:"
    echo "   https://taskfile.dev/installation/"
    exit 1
fi

# Test Go build
echo "1. Testing Go build..."
if task build; then
    echo "✅ Go build successful"
    task clean
else
    echo "❌ Go build failed"
    exit 1
fi

# Test Go tests
echo "2. Running Go tests..."
if task test; then
    echo "✅ Go tests passed"
else
    echo "❌ Go tests failed"
    exit 1
fi

# Test Docker build
echo "3. Testing Docker build..."
if task docker-build; then
    echo "✅ Docker build successful"
    docker rmi optiassign:latest
else
    echo "❌ Docker build failed"
    exit 1
fi

echo "🎉 All tests passed! OptiAssign is ready for development."
echo ""
echo "🔧 Available Task commands:"
echo "   task --list              # Show all available tasks"
echo "   task dev                 # Start development server"
echo "   task test-coverage       # Run tests with coverage"
echo "   task lint                # Run linter"