# DigitalOcean terraform NS8 cluster environment 

Terraform configuration for create set of Droplets to use as base for create a NS8 cluster.
The droplets will be created with DNS records and DigitalOcean firewall already configured,
but the NS8 system must be installed and configured manually.

## Variables

* `do_token`: DigitalOcean token.
* `sshkey`: DigitalOcean ssh public key to use (optional).
* `project`: DigitalOcean project where to create the droplets.
* `domain`: DigitalOcean domain where to create the DNS records.
* `leader_node`: The leader node (optional).
* `worker_nodes`: List of worker nodes (optional).

### The `leader_node` and `worker_nodes` variables

Each variable is a map and each item represents a cluster node.

- The item _key_ selects the OS type. Specify it followed by a number:

  * `dn` is for Debian 11
  * `fc` is for Fedora 34
  * `cs` is for CentOS Stream 9

- The item _value_ selects the droplet region. Refer to `doctl compute region list` output for
  a list of valid region codes.

The variable `leader_node` represents the leader node and the `worker_nodes` represents the worker nodes.

#### Centos Stream 9 nodes

Because CentOS Stream 9 is not present on DigitalOcean, a custom image must be created:

1. Upload the official [CentOS Cloud image](https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-20220121.1.x86_64.qcow2),
   following the [DigitalOcean documentation](https://docs.digitalocean.com/products/images/custom-images/how-to/upload/#image-requirements)
1. Set the name of the image to `CentOS-Stream-GenericCloud-9-20220121.1`

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

       terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"fc1":"sfo3"}'
       # -> dn1.leader.cluster0.dp.nethserver.net
       # -> fc1.worker.cluster0.dp.nethserver.net

5. Add another node to it:

       terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"fc1":"sfo3","fc2:"lon1"}'
       # -> fc2.worker.cluster0.dp.nethserver.net

6. Destroy the cluster

       terraform destroy

To work with multiple cluster instances just add more Terraform
workspaces. E.g.:

    terraform workspace new cluster1
    terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"fc5":"ams3","fc6":"sfo3","fc7":"sgp1"}'
    terraform workspace select cluster0
    terraform apply -var 'leader_node={"dn1":"ams3"}' -var 'worker_nodes={"dn1":"ams3","fc2":"sfo3"}'

## Shared firewall configuration

A shared firewall resource is applied to each node of `cluster0`. If a
different firewall configuration is required by some node, create another
firewall resource and assign it to that node.

## Default SSH keys pair

A pair of public and private key will be crated and installed on the cluster, for retrive the private key
and use it:

    terraform output -raw deploy-key  > key
