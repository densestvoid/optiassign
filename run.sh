#!/bin/bash

# OptiAssign Development Script

echo "Starting OptiAssign development environment..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "Please edit .env file with your Google OAuth credentials"
    exit 1
fi

# Start services with docker-compose
echo "Starting services..."
docker-compose up -d

# Wait for database to be ready
echo "Waiting for database to be ready..."
sleep 10

# Run database migrations
echo "Running database migrations..."
docker-compose exec -T app psql $DATABASE_URL -f migrations/001_initial_schema.sql

echo "Application is ready!"
echo "Access the application at: http://localhost:8080"
echo "To stop the application, run: docker-compose down"