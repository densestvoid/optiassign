#!/bin/bash

# OptiAssign Setup Script

echo "Setting up OptiAssign..."

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first:"
    echo "   https://golang.org/dl/"
    exit 1
fi

# Install tools if needed
echo "Installing Go tools..."
go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/pressly/goose/v3/cmd/goose@latest

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
    echo ""
    echo "After editing .env, run: task dev-setup"
    exit 1
fi

# Use Task for setup
echo "Running development setup with Task..."
go run github.com/go-task/task/v3/cmd/task@latest dev-setup

echo "✅ Setup complete!"
echo "🌐 Access the application at: http://localhost:8080"
echo ""
echo "🔧 Common Task commands:"
echo "   task dev              # Start development server"
echo "   task docker-compose-up # Start services"
echo "   task db-migrate       # Run migrations"
echo "   task db-status        # Check migration status"
echo "   task test             # Run tests"
echo "   task build            # Build application"
echo ""
echo "📝 To add pre-registered users, connect to the database and run:"
echo "   INSERT INTO users (google_id, email, name) VALUES ('', 'user@example.com', 'User Name');"