#!/bin/bash

# OptiAssign Production Deployment Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="optiassign"
AWS_REGION="us-west-2"
TERRAFORM_DIR="deployments/terraform"

echo -e "${GREEN}🚀 OptiAssign Production Deployment${NC}"
echo "=================================="

# Check prerequisites
echo -e "${YELLOW}📋 Checking prerequisites...${NC}"

# Check if required tools are installed
command -v aws >/dev/null 2>&1 || { echo -e "${RED}❌ AWS CLI is required but not installed.${NC}"; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo -e "${RED}❌ Terraform is required but not installed.${NC}"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo -e "${RED}❌ Docker is required but not installed.${NC}"; exit 1; }

echo -e "${GREEN}✅ All prerequisites found${NC}"

# Check AWS credentials
echo -e "${YELLOW}🔐 Checking AWS credentials...${NC}"
aws sts get-caller-identity >/dev/null 2>&1 || { echo -e "${RED}❌ AWS credentials not configured.${NC}"; exit 1; }
echo -e "${GREEN}✅ AWS credentials configured${NC}"

# Build and push Docker image
echo -e "${YELLOW}🐳 Building and pushing Docker image...${NC}"

# Get AWS account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REGISTRY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
ECR_REPOSITORY="${ECR_REGISTRY}/${APP_NAME}"

# Create ECR repository if it doesn't exist
aws ecr describe-repositories --repository-names ${APP_NAME} --region ${AWS_REGION} >/dev/null 2>&1 || {
    echo -e "${YELLOW}📦 Creating ECR repository...${NC}"
    aws ecr create-repository --repository-name ${APP_NAME} --region ${AWS_REGION}
}

# Login to ECR
echo -e "${YELLOW}🔑 Logging in to ECR...${NC}"
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR_REGISTRY}

# Build Docker image
echo -e "${YELLOW}🔨 Building Docker image...${NC}"
docker build -t ${APP_NAME}:latest .

# Tag and push image
echo -e "${YELLOW}📤 Pushing Docker image...${NC}"
docker tag ${APP_NAME}:latest ${ECR_REPOSITORY}:latest
docker push ${ECR_REPOSITORY}:latest

echo -e "${GREEN}✅ Docker image pushed successfully${NC}"

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
DATABASE_ENDPOINT=$(terraform output -raw database_endpoint)

echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo "=================================="
echo -e "${GREEN}Application URL: ${APPLICATION_URL}${NC}"
echo -e "${GREEN}Database Endpoint: ${DATABASE_ENDPOINT}${NC}"
echo ""
echo -e "${YELLOW}📝 Next steps:${NC}"
echo "1. Update your Google OAuth redirect URI to: ${APPLICATION_URL}/auth/google/callback"
echo "2. Run database migrations: go tool task db-migrate"
echo "3. Test the application at: ${APPLICATION_URL}"
echo ""
echo -e "${GREEN}✨ OptiAssign is now deployed and ready to use!${NC}"