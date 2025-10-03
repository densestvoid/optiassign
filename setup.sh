#!/bin/bash

# OptiAssign Setup Script

echo "Setting up OptiAssign..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your Google OAuth credentials before running the application"
    echo "   You need to:"
    echo "   1. Go to https://console.cloud.google.com/"
    echo "   2. Create a new project or select existing one"
    echo "   3. Enable Google+ API"
    echo "   4. Create OAuth 2.0 credentials"
    echo "   5. Add authorized redirect URI: http://localhost:8080/auth/google/callback"
    echo "   6. Update GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET in .env file"
    exit 1
fi

# Start services with docker-compose
echo "Starting services..."
docker-compose up -d

# Wait for database to be ready
echo "Waiting for database to be ready..."
sleep 10

# Run database migrations using Goose
echo "Running database migrations with Goose..."
docker-compose exec app go run cmd/migrate/main.go -command=up

echo "✅ Setup complete!"
echo "🌐 Access the application at: http://localhost:8080"
echo "🛑 To stop the application, run: docker-compose down"
echo ""
echo "📝 To add pre-registered users, connect to the database and run:"
echo "   INSERT INTO users (google_id, email, name) VALUES ('', 'user@example.com', 'User Name');"
echo ""
echo "🔧 Migration commands:"
echo "   go run cmd/migrate/main.go -command=up     # Run migrations"
echo "   go run cmd/migrate/main.go -command=down   # Rollback migrations"
echo "   go run cmd/migrate/main.go -command=status # Check migration status"