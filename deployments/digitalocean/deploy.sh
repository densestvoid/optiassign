#!/bin/bash

# OptiAssign DigitalOcean Deployment Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="optiassign"
TERRAFORM_DIR="deployments/digitalocean"

echo -e "${GREEN}🚀 OptiAssign DigitalOcean Deployment${NC}"
echo "======================================"

# Check prerequisites
echo -e "${YELLOW}📋 Checking prerequisites...${NC}"

# Check if required tools are installed
command -v doctl >/dev/null 2>&1 || { echo -e "${RED}❌ DigitalOcean CLI (doctl) is required but not installed.${NC}"; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo -e "${RED}❌ Terraform is required but not installed.${NC}"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo -e "${RED}❌ Docker is required but not installed.${NC}"; exit 1; }

echo -e "${GREEN}✅ All prerequisites found${NC}"

# Check DigitalOcean credentials
echo -e "${YELLOW}🔐 Checking DigitalOcean credentials...${NC}"
doctl auth list >/dev/null 2>&1 || { echo -e "${RED}❌ DigitalOcean credentials not configured. Run 'doctl auth init'${NC}"; exit 1; }
echo -e "${GREEN}✅ DigitalOcean credentials configured${NC}"

# Check if terraform.tfvars exists
if [ ! -f "${TERRAFORM_DIR}/terraform.tfvars" ]; then
    echo -e "${YELLOW}⚠️  terraform.tfvars not found. Please copy terraform.tfvars.example and configure your values.${NC}"
    echo -e "${BLUE}📝 Required variables:${NC}"
    echo "  - do_token: Your DigitalOcean API token"
    echo "  - db_password: Database password"
    echo "  - google_client_id: Google OAuth Client ID"
    echo "  - google_client_secret: Google OAuth Client Secret"
    echo "  - session_secret: Session secret key"
    echo "  - github_repo: Your GitHub repository (format: owner/repo)"
    exit 1
fi

# Deploy infrastructure with Terraform
echo -e "${YELLOW}🏗️  Deploying infrastructure with Terraform...${NC}"

cd ${TERRAFORM_DIR}

# Initialize Terraform
echo -e "${YELLOW}🔧 Initializing Terraform...${NC}"
terraform init

# Plan deployment
echo -e "${YELLOW}📋 Planning Terraform deployment...${NC}"
terraform plan -out=tfplan

# Apply deployment
echo -e "${YELLOW}🚀 Applying Terraform deployment...${NC}"
terraform apply tfplan

# Get outputs
echo -e "${YELLOW}📊 Getting deployment outputs...${NC}"
APPLICATION_URL=$(terraform output -raw application_url)
DATABASE_HOST=$(terraform output -raw database_host)
DATABASE_PORT=$(terraform output -raw database_port)
DATABASE_NAME=$(terraform output -raw database_name)
DATABASE_USER=$(terraform output -raw database_user)

echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo "======================================"
echo -e "${GREEN}Application URL: ${APPLICATION_URL}${NC}"
echo -e "${GREEN}Database Host: ${DATABASE_HOST}${NC}"
echo -e "${GREEN}Database Port: ${DATABASE_PORT}${NC}"
echo -e "${GREEN}Database Name: ${DATABASE_NAME}${NC}"
echo -e "${GREEN}Database User: ${DATABASE_USER}${NC}"
echo ""
echo -e "${YELLOW}📝 Next steps:${NC}"
echo "1. Update your Google OAuth redirect URI to: ${APPLICATION_URL}/auth/google/callback"
echo "2. Run database migrations: go tool task db-migrate"
echo "3. Test the application at: ${APPLICATION_URL}"
echo ""
echo -e "${BLUE}🔧 Useful commands:${NC}"
echo "  View app logs: doctl apps logs ${APP_NAME}-production"
echo "  Scale app: doctl apps update <app-id> --spec .do/app.yaml"
echo "  View database: doctl databases get <db-id>"
echo ""
echo -e "${GREEN}✨ OptiAssign is now deployed on DigitalOcean!${NC}"