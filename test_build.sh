#!/bin/bash

# Test build script for OptiAssign

echo "Testing OptiAssign build..."

# Test Go build
echo "1. Testing Go build..."
if go build -o main ./cmd/server; then
    echo "✅ Go build successful"
    rm -f main
else
    echo "❌ Go build failed"
    exit 1
fi

# Test Go tests
echo "2. Running Go tests..."
if go test ./domain/...; then
    echo "✅ Go tests passed"
else
    echo "❌ Go tests failed"
    exit 1
fi

# Test Docker build
echo "3. Testing Docker build..."
if docker build -t optiassign-test .; then
    echo "✅ Docker build successful"
    docker rmi optiassign-test
else
    echo "❌ Docker build failed"
    exit 1
fi

echo "🎉 All tests passed! OptiAssign is ready for development."