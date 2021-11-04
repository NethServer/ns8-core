
variable "images" {
  description = "Map host name code to OS image"
  default = {
    "fc" = "fedora-34-x64",
    "dn" = "debian-11-x64",
    "ub" = "ubuntu-21-04-x64"
  }
}

resource "digitalocean_droplet" "leader" {
  for_each           = var.leader_nodes
  image              = var.images[substr(each.key, 0, 2)]
  name               = format("%s.leader.%s.%s", each.key, terraform.workspace ,var.domain)
  region             = each.value
  size               = "s-1vcpu-1gb-intel"
  ipv6               = substr(each.key, 0, 2) != "dn"
  private_networking = true
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
  tags = [digitalocean_tag.cluster.name]
}

resource "digitalocean_droplet" "worker" {
  for_each           = var.worker_nodes
  image              = var.images[substr(each.key, 0, 2)]
  name               = format("%s.worker.%s.%s", each.key, terraform.workspace ,var.domain)
  region             = each.value
  size               = "s-1vcpu-1gb-intel"
  ipv6               = substr(each.key, 0, 2) != "dn"
  private_networking = true
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
  tags = [digitalocean_tag.cluster.name]
}

resource "digitalocean_project_resources" "leader_nodes" {
  project   = data.digitalocean_project.default.id
  resources = [for hpx, rgn in var.leader_nodes : digitalocean_droplet.leader[hpx].urn]
}

resource "digitalocean_project_resources" "worker_nodes" {
  project   = data.digitalocean_project.default.id
  resources = [for hpx, rgn in var.worker_nodes : digitalocean_droplet.worker[hpx].urn]
}


resource "digitalocean_record" "leader" {
  for_each = var.leader_nodes
  type     = "A"
  domain   = data.digitalocean_domain.default.name
  name     = format("%s.leader.%s", each.key, terraform.workspace)
  value    = digitalocean_droplet.leader[each.key].ipv4_address
  ttl      = 300
}

resource "digitalocean_record" "worker" {
  for_each = var.worker_nodes
  type     = "A"
  domain   = data.digitalocean_domain.default.name
  value    = digitalocean_droplet.worker[each.key].ipv4_address
  name     = format("%s.worker.%s", each.key, terraform.workspace)
  ttl      = 300
}

resource "digitalocean_tag" "cluster" {
  name = format("%s-%s", var.project, terraform.workspace)
}
