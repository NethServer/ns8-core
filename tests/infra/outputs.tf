output "leader_node" {
  value = [for hpx, _ in digitalocean_droplet.leader : digitalocean_droplet.leader[hpx].name]
}

output "worker_nodes" {
  value = [for hpx, _ in digitalocean_droplet.worker : digitalocean_droplet.worker[hpx].name]
}

output "deploy-key" {
  value = tls_private_key.deploy.private_key_pem
  sensitive = true
}
