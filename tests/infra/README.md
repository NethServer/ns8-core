# Personal Terraform NS8 cluster environment

1. Create a `davidep.auto.tfvars` file, like the following:

       sshkey   = "davidep"
       do_token = "secret"
       project  = "davidep"
       domain   = "dp.nethserver.net"

2. Create and select a new workspace `cluster0`

       terraform workspace create cluster0

3. Create two nodes for `cluster0`

       terraform apply -var 'nodes={"dn1":"ams3","fc2":"sfo3"}'
       # -> dn1.dp.nethserver.net
       # -> fc2.dp.nethserver.net

4. Add another node to it:

       terraform apply -var 'nodes={"dn1":"ams3","fc2":"sfo3","fc3":"lon1"}'
       # -> fc3.dp.nethserver.net

5. Destroy the cluster

       terraform destroy -var 'nodes={}'

To work with multiple cluster instances just add more Terraform
workspaces. E.g.:

    terraform workspace new cluster1
    terraform apply -var 'nodes={"fc5":"ams3","fc6":"sfo3","fc7":"sgp1"}'
    terraform workspace select cluster0
    terraform destroy -var 'nodes={"dn1":"ams3","fc2":"sfo3"}'

## The `nodes` variable

The `nodes` variable is a map. Each item represents a cluster node.

- The item _key_ selects the OS type. Specify it followed by a number:

  * `dn` is for Debian
  * `fc` is for Fedora
  * `ub` is for Ubuntu

- The item _value_ selects the VPS region. Refer to `doctl compute region list` output for
  a list of valid region codes.

## Shared firewall configuration

A shared firewall resource is applied to each node of `cluster0`. If a
different firewall configuration is required by some node, create another
firewall resource and assign it to that node.
