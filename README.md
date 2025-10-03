# OptiAssign - Minimal Viable Solution

A web application for group-based, randomized, prioritized, snaking-draft item assignment.

## ✅ Phase 1 Complete: Setup & Auth

**Current Status**: Phase 1 implementation is complete with working Google SSO authentication, database connection, and minimal layout.

### Features Implemented

- ✅ Google SSO Authentication (OAuth 2.0)
- ✅ PostgreSQL database with complete schema
- ✅ HTMX-based frontend with Bootstrap 5
- ✅ Go templating system
- ✅ Docker containerization
- ✅ Unit tests for core authentication logic
- ✅ Secure session management

## Setup

### Prerequisites

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

## Next Steps (Phase 2)

The next phase will implement:
- Group creation and management
- Item CRUD operations
- Participant management
- HTMX forms with server-side validation