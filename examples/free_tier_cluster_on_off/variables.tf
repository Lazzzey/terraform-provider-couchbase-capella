variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

variable "organization_id" {
  description = "Capella Organization ID"
}

variable "project_id" {
  description = "Project Name for Project Created via Terraform"
}

variable "cluster_id" {
  description = "Capella Cluster ID"
}

variable "state" {
  description = "Provides the state cluster to be switched to- on or off"
}
