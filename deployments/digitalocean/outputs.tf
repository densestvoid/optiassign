# OptiAssign DigitalOcean Outputs

output "application_url" {
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

# Note: Using GitHub Container Registry (ghcr.io) - no additional output needed

output "vpc_id" {
  description = "VPC ID"
  value       = digitalocean_vpc.main.id
}

output "app_id" {
  description = "App Platform application ID"
  value       = digitalocean_app.main.id
}

output "database_cluster_id" {
  description = "Database cluster ID"
  value       = digitalocean_database_cluster.main.id
}

output "firewall_id" {
  description = "Firewall ID"
  value       = digitalocean_firewall.main.id
}