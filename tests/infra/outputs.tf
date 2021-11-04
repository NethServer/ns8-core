output "leader_node" {
  value = [for hpx, _ in digitalocean_droplet.leader : digitalocean_droplet.leader[hpx].name]
}

output "worker_nodes" {
  value = [for hpx, _ in digitalocean_droplet.worker : digitalocean_droplet.worker[hpx].name]
}
