# samba

The Samba module allows to install one or more AD domains in a NS8
cluster. An AD domain (or realm, in Kerberos terms), can have one or more
domain controllers (DCs). Provisioning a new domain and joining an
additional domain controller to an existing domain are both implemented by
the same action: `provision-domain`.

A Samba DC requires a dedicated IP address. Two DCs cannot share the same
IP.

Do not assign the IP address of a untrusted network! The IP address
assigned to Samba DC exposes internal services that are not designed to be
exposed publicly. The IP address can be added to the cluster VPN routes,
so DCs can see each other and perform replication.

The LDAP service of Samba DCs does not allow clear-text LDAP binds: enable
TLS or use Kerberos/GSSAPI authentication. The default node Ldapproxy
instance is configured to enable TLS automatically.

To run the example commands of the following sections initialize a shell
environment with the following commands:

    pip install httpie
    export TOKEN=$(http :8080/api/login username=admin password=Nethesis,1234 | jq -r .token)

## Provision

Provision a new domain:

    http :8080/api/module/samba1/tasks "Authorization: Bearer $TOKEN" <<EOF
    {
        "action":"provision-domain",
        "data":{
            "adminuser":"administrator",
            "adminpass":"Nethesis,1234",
            "realm":"ad.$(hostname -d)",
            "nbdomain":"ad",
            "ipaddress":"10.133.0.2"
        }
    }
    EOF

Further Samba instances for the same `realm` are **joined** to the existing domain.
The command is similar.

    http :8080/api/module/samba2/tasks "Authorization: Bearer $TOKEN" <<EOF
    {
        "action":"provision-domain",
        "data":{
            "adminuser":"administrator",
            "adminpass":"Nethesis,1234",
            "realm":"ad.$(hostname -d)",
            "ipaddress":"10.124.0.2"
        }
    }
    EOF

Note that the `nbdomain` attribute is no longer required to provision
additional DCs.

## IP routing for the AD domain

In multiple DCs domains, DCs must connect each other to join the domain
and for data replication. If a DC cannot contact other DCs the provision
procedure fails.

If the existing network infrastructure does not provide the necessary IP
routes, it is possible to use the cluster VPN to route IP traffic among
DCs. For example this command enable the cluster VPN routes to join a second
DC:

    http :8080/api/module/samba2/tasks "Authorization: Bearer $TOKEN" <<EOF
    {
        "action":"set-routes",
        "data": {
            "realm": "ad.$(hostname -d)",
            "ipaddress": "10.124.0.2",
            "use_cluster_vpn": true
        }
    }
    EOF

## Configure Ldapproxy

When a Samba module is provisioned for the first time on a node, it checks
if the node default Ldapproxy instance is not bound to another account
provider. If so, Samba invokes the `set-backend` on Ldapproxy
automatically.

This command manually binds `ldapproxy1` with the `samba2` local account
provider LDAP backend. TLS is required by Samba (clear-text LDAP binds are
not allowed).

    # retrieve the ldapservice password
    redis-cli HGET module/samba1/environment SVCPASS
    # output: "Random,1234"
    http :8080/api/module/ldapproxy1/tasks "Authorization: Bearer $TOKEN" <<EOF
    {
        "action":"set-backend",
        "data": {
            "backend": "samba1",
            "schema": "ad",
            "host": "127.0.0.1",
            "port": 636,
            "tls": true,
            "tls_verify": false,
            "bind_dn": "AD\\ldapservice",
            "bind_password:": "Random,1234"
        }
    }
    EOF
