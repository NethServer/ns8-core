# ldapproxy

The Ldapproxy module is a communication broker between an account provider
LDAP backend (like Samba DC, OpenLDAP...) and a consumer module that wants
to authenticate its users against the LDAP service (e.g. Nextcloud).

## Connect to the LDAP service proxy

The consumer module runs on the same node where the Ldapproxy instance is
running. A default Ldapproxy instance is installed as a core module on
each cluster node.

1. Retrieve the default instance ID

       HGET node/${NODE_ID}/default_instance/ldapproxy (=> ldapproxy1)

2. Retrieve the TCP port

       HGET module/ldapproxy1/environment TCP_PORT (=> 20000)

If the module LDAP client runs in a Podman container with a `host`
network, just point the LDAP client to `ldap://127.0.0.1:20000`.

If the module LDAP client runs in a Podman container with a `private`
network, add the following arguments to the `podman run` command:

    --network=slirp4netns:allow_host_loopback=true --add-host=accountprovider:10.0.2.2

Point the client to `ldap://accountprovider:20000`.

## Set the proxy LDAP backend

Modules with role `accountprovider` can invoke the `set-backend` action on
Ldapproxy instances.

The `set-backend` action configures and restarts the L4 proxy service.

    # pip install httpie
    export TOKEN=$(http :8080/api/login username=admin password=Nethesis,1234 | jq -r .token)
    http :8080/api/module/ldapproxy1/tasks "Authorization: Bearer $TOKEN" <<EOF
    {
        "action":"set-backend", 
        "data": {
            "backend": "samba1",
            "schema": "ad",
            "host": "127.0.0.1",
            "port": 636,
            "tls": true,
            "tls_verify": false
        }
    }
    EOF
