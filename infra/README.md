# DigitalOcean terraform NS8 cluster environment 

Terraform configuration for create set of Droplets to use as base for create a NS8 cluster.
The droplets will be created with DNS records and DigitalOcean firewall already configured,
but the NS8 system must be installed and configured manually.

## Variables

* `do_token`: DigitalOcean token.
* `sshkey`: DigitalOcean ssh public key to use (optional).
   If your private key is password-protected, do not set this variable and use the SSH automatically generated
   by terraform.
* `project`: DigitalOcean project where to create the droplets.
* `domain`: DigitalOcean domain where to create the DNS records.
* `leader_node`: The leader node (optional).
* `worker_nodes`: List of worker nodes (optional).

### The `leader_node` and `worker_nodes` variables

Each variable is a map and each item represents a cluster node.

- The item _key_ selects the OS type. Specify it followed by a number:

  * `dn` is for Debian 11
  * `cs` is for CentOS Stream 9

- The item _value_ selects the droplet region. Refer to `doctl compute region list` output for
  a list of valid region codes.

The variable `leader_node` represents the leader node and the `worker_nodes` represents the worker nodes.

## Examples

Download and install [Terraform](https://www.terraform.io/downloads), then follow below steps.

1. Create a `configs.auto.tfvars` file, like the following:

       sshkey   = "davidep"
       do_token = "secret"
       project  = "davidep"
       domain   = "dp.nethserver.net"

2. Install required plugins

       terraform init

3. Create and select a new workspace `cluster0`

       terraform workspace new cluster0

4. Create two nodes for `cluster0`

       terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"cs1":"sfo3"}'
       # -> dn1.leader.cluster0.dp.nethserver.net
       # -> cs1.worker.cluster0.dp.nethserver.net

5. Add another node to it:

       terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"cs1":"sfo3","cs2:"lon1"}'
       # -> cs2.worker.cluster0.dp.nethserver.net

6. Destroy the cluster

       terraform destroy

To work with multiple cluster instances just add more Terraform
workspaces. E.g.:

    terraform workspace new cluster1
    terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"dn5":"ams3","dn6":"sfo3","dn7":"sgp1"}'
    terraform workspace select cluster0
    terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"dn1":"ams3","dn2":"sfo3"}'

## Shared firewall configuration

A shared firewall resource is applied to each node of `cluster0`. If a
different firewall configuration is required by some node, create another
firewall resource and assign it to that node.

## Default SSH keys pair

A pair of public and private key will be crated and installed on the cluster, for retrive the private key
and use it:

    terraform output -raw deploy-key  > key
    chmod 0600 key
