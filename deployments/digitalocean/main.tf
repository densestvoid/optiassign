# OptiAssign DigitalOcean Infrastructure
terraform {
  required_version = ">= 1.0"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = var.do_token
}

# Variables
variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "app_name" {
  description = "Application name"
  type        = string
  default     = "optiassign"
}

variable "region" {
  description = "DigitalOcean region"
  type        = string
  default     = "nyc3"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "google_client_id" {
  description = "Google OAuth Client ID"
  type        = string
  sensitive   = true
}

variable "google_client_secret" {
  description = "Google OAuth Client Secret"
  type        = string
  sensitive   = true
}

variable "session_secret" {
  description = "Session secret for authentication"
  type        = string
  sensitive   = true
}

# VPC
resource "digitalocean_vpc" "main" {
  name     = "${var.app_name}-${var.environment}-vpc"
  region   = var.region
  ip_range = "10.10.0.0/16"
}

# SSH Key (optional - for debugging)
resource "digitalocean_ssh_key" "main" {
  name       = "${var.app_name}-${var.environment}-key"
  public_key = file("~/.ssh/id_rsa.pub")
}

# Managed Database
resource "digitalocean_database_cluster" "main" {
  name       = "${var.app_name}-${var.environment}-db"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = var.region
  node_count = 1
  vpc_uuid   = digitalocean_vpc.main.id

  tags = [
    "${var.app_name}",
    var.environment,
    "database"
  ]
}

# Database User
resource "digitalocean_database_user" "main" {
  cluster_id = digitalocean_database_cluster.main.id
  name       = "optiassign"
}

# Database Database
resource "digitalocean_database_db" "main" {
  cluster_id = digitalocean_database_cluster.main.id
  name       = "optiassign"
}

# Note: Using GitHub Container Registry (ghcr.io) instead of DigitalOcean registry
# This saves costs and simplifies the deployment process

# App Platform App
resource "digitalocean_app" "main" {
  spec {
    name   = "${var.app_name}-${var.environment}"
    region = var.region

    # Static site for frontend (optional)
    static_site {
      name    = "frontend"
      source_dir = "/"
      build_command = "echo 'Static site'"
      output_dir = "/"
    }

    # Service for backend
    service {
      name               = "backend"
      source_dir         = "/"
      github {
        repo           = var.github_repo
        branch         = var.github_branch
        deploy_on_push = true
      }

      # Build configuration
      build_command = "go build -o main ./cmd/server"
      run_command   = "./main"

      # Environment variables
      env {
        key   = "PORT"
        value = "8080"
        type  = "GENERAL"
      }

      env {
        key   = "DB_HOST"
        value = digitalocean_database_cluster.main.host
        type  = "GENERAL"
      }

      env {
        key   = "DB_PORT"
        value = tostring(digitalocean_database_cluster.main.port)
        type  = "GENERAL"
      }

      env {
        key   = "DB_NAME"
        value = digitalocean_database_db.main.name
        type  = "GENERAL"
      }

      env {
        key   = "DB_USER"
        value = digitalocean_database_user.main.name
        type  = "GENERAL"
      }

      env {
        key   = "DB_PASSWORD"
        value = var.db_password
        type  = "SECRET"
      }

      env {
        key   = "GOOGLE_CLIENT_ID"
        value = var.google_client_id
        type  = "SECRET"
      }

      env {
        key   = "GOOGLE_CLIENT_SECRET"
        value = var.google_client_secret
        type  = "SECRET"
      }

      env {
        key   = "SESSION_SECRET"
        value = var.session_secret
        type  = "SECRET"
      }

      env {
        key   = "EMAIL_LOG_TO_CONSOLE"
        value = "true"
        type  = "GENERAL"
      }

      # Instance configuration
      instance_count    = 1
      instance_size_slug = "basic-xxs"

      # Health check
      health_check {
        http_path             = "/"
        initial_delay_seconds = 10
        period_seconds        = 10
        timeout_seconds       = 5
        success_threshold     = 1
        failure_threshold     = 3
      }

      # Routes
      routes {
        path = "/"
      }

      routes {
        path = "/auth"
      }

      routes {
        path = "/groups"
      }

      routes {
        path = "/participant"
      }
    }

    # Database connection
    database {
      name = "optiassign-db"
      engine = "PG"
      cluster_name = digitalocean_database_cluster.main.name
      db_name = digitalocean_database_db.main.name
      db_user = digitalocean_database_user.main.name
    }
  }
}

# Firewall
resource "digitalocean_firewall" "main" {
  name = "${var.app_name}-${var.environment}-firewall"

  droplet_ids = []

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  tags = [
    "${var.app_name}",
    var.environment
  ]
}

# Domain (optional)
resource "digitalocean_domain" "main" {
  count = var.domain_name != "" ? 1 : 0
  name  = var.domain_name
}

# Domain record
resource "digitalocean_record" "main" {
  count  = var.domain_name != "" ? 1 : 0
  domain = digitalocean_domain.main[0].name
  type   = "CNAME"
  name   = "@"
  value  = digitalocean_app.main.default_ingress
  ttl    = 300
}

# Outputs
output "app_url" {
  description = "URL of the deployed application"
  value       = digitalocean_app.main.default_ingress
}

output "database_host" {
  description = "Database host"
  value       = digitalocean_database_cluster.main.host
  sensitive   = true
}

output "database_port" {
  description = "Database port"
  value       = digitalocean_database_cluster.main.port
}

output "database_name" {
  description = "Database name"
  value       = digitalocean_database_db.main.name
}

output "database_user" {
  description = "Database user"
  value       = digitalocean_database_user.main.name
}

output "container_registry" {
  description = "Container registry endpoint"
  value       = digitalocean_container_registry.main.server_url
}

output "vpc_id" {
  description = "VPC ID"
  value       = digitalocean_vpc.main.id
}