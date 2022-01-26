variable "project" {
  description = "Project name"
}

variable "domain" {
  description = "Personal DNS domain suffix"
  default     = "nethserver.net"
}

variable "leader_node" {
  description = "Leader node"
  type        = map(string)
  default = {}
  validation {
    condition = length(var.leader_node) <= 1
    error_message = "The leader node must be one or zero."
  }
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

