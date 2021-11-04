terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.0.0"
    }
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.3.0"
    }
  }
}

variable "sshkey" {
  description = "DigitalOcean SSH key name"
}

provider "digitalocean" {
  token = var.do_token
}

variable "do_token" {
  description = "DigitalOcean API token"
}

data "digitalocean_ssh_key" "terraform" {
  name = var.sshkey
}
