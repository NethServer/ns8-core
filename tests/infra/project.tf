variable "project" {
  description = "Project name"
}

variable "domain" {
  description = "Personal DNS domain suffix"
  default     = "nethserver.net"
}

variable "leader_nodes" {
  description = "List of leader nodes"
  type        = map(string)
  default = {}
}

variable "worker_nodes" {
  description = "List of worker nodes"
  type        = map(string)
  default = {}
}

data "digitalocean_project" "default" {
  name = var.project
}

data "digitalocean_domain" "default" {
  name = var.domain
}

