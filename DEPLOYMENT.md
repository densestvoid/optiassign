# OptiAssign DigitalOcean Deployment Guide

This guide covers deploying OptiAssign to DigitalOcean using Terraform and App Platform.

## Prerequisites

### Required Tools
- [DigitalOcean CLI (doctl)](https://docs.digitalocean.com/reference/doctl/how-to/install/)
- [Terraform](https://www.terraform.io/downloads)
- [Docker](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/downloads)

### Required Accounts
- DigitalOcean account with API token
- GitHub account with repository access
- Google Cloud Console account for OAuth

## Quick Deployment

### 1. Clone and Setup
```bash
git clone <your-repo-url>
cd optiassign
```

### 2. Configure DigitalOcean
```bash
# Install doctl
# macOS
brew install doctl

# Linux
snap install doctl

# Authenticate
doctl auth init
```

### 3. Configure Terraform Variables
```bash
cd deployments/digitalocean
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your values:
```hcl
# DigitalOcean Configuration
do_token = "your-digitalocean-api-token"
region = "nyc3"
environment = "production"
app_name = "optiassign"

# Database Configuration
db_password = "your-secure-database-password"

# Google OAuth Configuration
google_client_id = "your-google-client-id"
google_client_secret = "your-google-client-secret"

# Session Configuration
session_secret = "your-session-secret-key"

# GitHub Configuration
github_repo = "your-username/optiassign"
github_branch = "main"
```

### 4. Deploy
```bash
./deploy.sh
```

## Manual Deployment

### 1. Create DigitalOcean Resources

#### Database
```bash
# Create managed database
doctl databases create optiassign-db \
  --engine pg \
  --version 15 \
  --size db-s-1vcpu-1gb \
  --region nyc3
```

#### Container Registry
```bash
# Create container registry
doctl registry create optiassign-registry
```

### 2. Build and Push Docker Image
```bash
# Login to registry
doctl registry login

# Build image
docker build -t optiassign:latest .

# Tag and push
docker tag optiassign:latest registry.digitalocean.com/optiassign-registry/optiassign:latest
docker push registry.digitalocean.com/optiassign-registry/optiassign:latest
```

### 3. Deploy App
```bash
# Create app from .do/app.yaml
doctl apps create --spec .do/app.yaml
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | Application port | Yes |
| `DB_HOST` | Database host | Yes |
| `DB_PORT` | Database port | Yes |
| `DB_NAME` | Database name | Yes |
| `DB_USER` | Database user | Yes |
| `DB_PASSWORD` | Database password | Yes |
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | Yes |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | Yes |
| `SESSION_SECRET` | Session secret key | Yes |
| `EMAIL_LOG_TO_CONSOLE` | Email logging | No |

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URI: `https://your-app-url.ondigitalocean.app/auth/google/callback`
6. Copy Client ID and Secret to your configuration

### Database Setup

1. Run migrations after deployment:
```bash
# Connect to your app
doctl apps ssh <app-id>

# Run migrations
go tool task db-migrate
```

## Monitoring and Management

### View Logs
```bash
# App logs
doctl apps logs <app-id>

# Database logs
doctl databases logs <db-id>
```

### Scale Application
```bash
# Update app spec
doctl apps update <app-id> --spec .do/app.yaml
```

### Database Management
```bash
# Get database info
doctl databases get <db-id>

# Create database user
doctl databases users create <db-id> --username optiassign
```

## Troubleshooting

### Common Issues

1. **App won't start**
   - Check environment variables
   - Verify database connection
   - Check build logs

2. **Database connection failed**
   - Verify database is running
   - Check firewall rules
   - Verify credentials

3. **OAuth not working**
   - Check redirect URI configuration
   - Verify Google OAuth setup
   - Check environment variables

### Debug Commands
```bash
# Check app status
doctl apps get <app-id>

# View app spec
doctl apps spec get <app-id>

# Check database status
doctl databases get <db-id>

# View app logs
doctl apps logs <app-id> --follow
```

## Cost Estimation

### DigitalOcean Pricing (Monthly)

| Resource | Size | Cost |
|----------|------|------|
| App Platform | Basic XXS | $5 |
| Managed Database | 1GB RAM | $15 |
| Container Registry | 5GB | $5 |
| **Total** | | **$25/month** |

### Scaling Options

- **App Platform**: Scale from Basic XXS to Professional plans
- **Database**: Upgrade to larger instances
- **Storage**: Add more registry storage as needed

## Security

### Best Practices
- Use strong passwords for database
- Rotate session secrets regularly
- Enable VPC for network isolation
- Use firewall rules to restrict access
- Monitor logs for suspicious activity

### Environment Security
- Store secrets in DigitalOcean App Platform environment variables
- Use different credentials for different environments
- Regularly update dependencies
- Monitor for security vulnerabilities

## Support

For issues with:
- **DigitalOcean**: Check [DigitalOcean Documentation](https://docs.digitalocean.com/)
- **Terraform**: Check [Terraform Documentation](https://www.terraform.io/docs/)
- **Application**: Check application logs and GitHub issues