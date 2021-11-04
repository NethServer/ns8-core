variable "project" {
  description = "Project name"
}

variable "domain" {
  description = "Personal DNS domain suffix"
  default     = "nethserver.net"
}

variable "nodes" {
  description = "Host name for the VPS"
  type        = map(string)
}

data "digitalocean_project" "default" {
  name = var.project
}

data "digitalocean_domain" "default" {
  name = var.domain
}

