# OptiAssign - Minimal Viable Solution

A web application for group-based, randomized, prioritized, snaking-draft item assignment.

## ✅ Phase 4 Complete: Finalization

**Current Status**: Phase 4 implementation is complete with assignment results display, production Docker deployment, and Terraform infrastructure.

### Features Implemented

- ✅ Google SSO Authentication (OAuth 2.0)
- ✅ PostgreSQL database with complete schema
- ✅ HTMX-based frontend with Bootstrap 5
- ✅ Go templating system
- ✅ Docker containerization
- ✅ Unit tests for core authentication logic
- ✅ Secure session management
- ✅ Modern configuration management
- ✅ Task management system for repository operations
- ✅ Goose migration management with CLI tools
- ✅ Enhanced database connection pooling
- ✅ Group creation and management
- ✅ Item CRUD operations
- ✅ Participant management
- ✅ HTMX forms with server-side validation
- ✅ Complete repository layer
- ✅ Business logic services
- ✅ Tokenized participant access
- ✅ Priority list submission forms
- ✅ Snaking draft assignment algorithm
- ✅ Email invitation system
- ✅ Assignment execution and results
- ✅ Automatic assignment execution when last participant submits
- ✅ Background task manager for assignment processing
- ✅ Real-time assignment status tracking
- ✅ Personalized email notifications with assigned items
- ✅ Detailed item information in assignment emails
- ✅ Enhanced email templates with item descriptions
- ✅ **NEW**: Assignment results display pages
- ✅ **NEW**: Production-ready Docker deployment
- ✅ **NEW**: Complete Terraform infrastructure
- ✅ **NEW**: Automated deployment scripts

## Setup

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Google OAuth credentials

### Quick Start

1. **Setup environment**:
   ```bash
   ./setup.sh
   ```

2. **Test the build**:
   ```bash
   ./test_build.sh
   ```

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
6. Update `.env` file with your credentials

### Manual Setup

1. Copy `.env.example` to `.env` and fill in your Google OAuth credentials
2. Start services: `docker-compose up -d`
3. Run migrations: `docker-compose exec app psql $DATABASE_URL -f migrations/001_initial_schema.sql`
4. Access: http://localhost:8080

## Development

### Project Structure

```
/cmd/server         # Main entry point
/domain             # Core business logic (with tests)
/api                # HTTP handlers and middleware
/db                 # Database connection and queries
/web                # HTML templates
/migrations         # SQL files for database schema
/deployments        # Terraform files (placeholder)
```

### Adding Pre-registered Users

To add pre-registered users, insert them directly into the database:

```sql
INSERT INTO users (google_id, email, name) VALUES ('', 'user@example.com', 'User Name');
```

## API Endpoints (Phase 1)

- `GET /` - Home page (redirects to dashboard if authenticated)
- `GET /login` - Login page
- `GET /auth/google/login` - Initiate Google OAuth
- `GET /auth/google/callback` - Google OAuth callback
- `GET /logout` - Logout
- `GET /dashboard` - User dashboard (protected)

## Technology Stack

- **Backend**: Go with Chi router
- **Database**: PostgreSQL
- **Frontend**: HTMX + Bootstrap 5
- **Templating**: Go html/template
- **Authentication**: Google OAuth 2.0
- **Containerization**: Docker

## New Tools & Improvements

### Configuration Management
- Modern configuration system with validation
- Environment variable handling with defaults
- Type-safe configuration loading

### Task Management
- Taskfile.yml for common repository commands
- go-tasks (Task) for development workflow
- Standardized build, test, and deployment tasks

### Migration Management
- Goose-based migration system
- CLI tools for migration management
- Rollback and status checking capabilities

### Task Commands (using modern Go tools pattern)
```bash
# First, install tools
go get -tool github.com/go-task/task/v3/cmd/task@latest
go get -tool github.com/pressly/goose/v3/cmd/goose@latest

# Development
go tool task dev              # Start development server
go tool task build            # Build application
go tool task test             # Run tests
go tool task test-coverage    # Run tests with coverage
go tool task lint             # Run linter

# Database
go tool goose -dir migrations up     # Run migrations
go tool goose -dir migrations down   # Rollback migrations
go tool goose -dir migrations status # Check migration status

# Docker
go tool task docker-build     # Build Docker image
go tool task docker-run       # Run in Docker
go tool task docker-compose-up # Start services
go tool task docker-compose-down # Stop services

# Setup
go tool task setup            # Initial setup
go tool task dev-setup        # Complete development setup
go tool task install-deps    # Install dependencies

# Cleanup
go tool task clean            # Clean build artifacts
go tool task clean-docker     # Clean Docker resources
```

## API Endpoints (Phase 3)

### Group Management
- `GET /groups` - List user's groups
- `GET /groups/new` - Create group form
- `POST /groups` - Create new group (HTMX)
- `GET /groups/{id}` - View group details
- `POST /groups/{id}/execute` - Execute assignment algorithm

### Participant Access (Token-based)
- `GET /participant/{token}` - Participant access page
- `GET /participant/{token}/priorities` - Priority submission form
- `POST /participant/{token}/priorities` - Submit priorities
- `GET /participant/{token}/status` - Assignment status (JSON)
- `GET /participant/{token}/results` - Participant assignment results

### Assignment Results
- `GET /groups/{id}/results` - Full assignment results (owner only)
- `GET /participant/{token}/results` - Participant's assigned items

## Production Deployment

### Quick Start
```bash
# Deploy to AWS with Terraform
./deployments/deploy.sh
```

### Manual Deployment
```bash
# Build Docker image
docker build -t optiassign:latest .

# Run with docker-compose
docker-compose up -d

# Apply database migrations
go tool task db-migrate
```

### Infrastructure
- **AWS ECS Fargate** for containerized application
- **Application Load Balancer** for traffic distribution
- **RDS PostgreSQL** for database
- **CloudWatch** for logging and monitoring
- **VPC** with public/private subnets for security

## Project Complete! 🎉

**OptiAssign MVS** is now a fully functional, production-ready application with:
- Complete assignment workflow
- Automatic execution
- Personalized notifications
- Production infrastructure
- Scalable deployment