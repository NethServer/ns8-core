locals {
  //Map host name code to OS image
  images = {
    "dn" = "debian-11-x64",
    "cs" = data.digitalocean_image.centos.id
  }
}

data "digitalocean_image" "centos" {
  name = "CentOS-Stream-GenericCloud-9-20220121.1"
}

resource "digitalocean_droplet" "leader" {
  for_each           = var.leader_node
  image              = local.images[substr(each.key, 0, 2)]
  name               = format("%s.leader.%s.%s", each.key, terraform.workspace ,var.domain)
  region             = each.value
  size               = "s-1vcpu-1gb-intel"
  private_networking = true
  ssh_keys = compact([try(data.digitalocean_ssh_key.terraform[0].id, ""),
    digitalocean_ssh_key.deploy.id])

  tags = [digitalocean_tag.cluster.name]
}

resource "digitalocean_droplet" "worker" {
  for_each           = var.worker_nodes
  image              = local.images[substr(each.key, 0, 2)]
  name               = format("%s.worker.%s.%s", each.key, terraform.workspace ,var.domain)
  region             = each.value
  size               = "s-1vcpu-1gb-intel"
  private_networking = true
  ssh_keys = compact([try(data.digitalocean_ssh_key.terraform[0].id, ""),
    digitalocean_ssh_key.deploy.id])
  tags = [digitalocean_tag.cluster.name]
}

resource "tls_private_key" "deploy" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "digitalocean_ssh_key" "deploy" {
  name       = format("%s.%s-deploy", terraform.workspace ,var.domain)
  public_key = tls_private_key.deploy.public_key_openssh
}


resource "digitalocean_project_resources" "leader_node" {
  project   = data.digitalocean_project.default.id
  resources = [for hpx, rgn in var.leader_node : digitalocean_droplet.leader[hpx].urn]
}

resource "digitalocean_project_resources" "worker_nodes" {
  project   = data.digitalocean_project.default.id
  resources = [for hpx, rgn in var.worker_nodes : digitalocean_droplet.worker[hpx].urn]
}


resource "digitalocean_record" "leader" {
  for_each = var.leader_node
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
