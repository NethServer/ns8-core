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


## Provision

Provision a new domain:

    api-cli run provision-domain --agent module/samba1 --data - <<EOF
    {
        "adminuser":"administrator",
        "adminpass":"Nethesis,1234",
        "realm":"ad.$(hostname -d)",
        "nbdomain":"ad",
        "ipaddress":"10.133.0.2"
    }
    EOF

Further Samba instances for the same `realm` are **joined** to the existing domain.
The command is similar.

    api-cli run provision-domain --agent module/samba2 --data - <<EOF
    {
        "adminuser":"administrator",
        "adminpass":"Nethesis,1234",
        "realm":"ad.$(hostname -d)",
        "ipaddress":"10.124.0.2"
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

    api-cli run set-routes --agent module/samba2 --data - <<EOF
    {
        "realm": "ad.$(hostname -d)",
        "ipaddress": "10.124.0.2",
        "use_cluster_vpn": true
    }
    EOF

## Configure Ldapproxy

When a Samba module is provisioned for the first time on a node, it checks
if the node default Ldapproxy instance is not bound to another account
provider. If so, Samba invokes the `set-backend` on Ldapproxy
automatically.

This command manually binds `ldapproxy3` with the `samba2` local account
provider LDAP backend. TLS is required by Samba (clear-text LDAP binds are
not allowed), however `tls_verify` is disabled because Samba has a
self-signed certificate.

    # retrieve the ldapservice password
    redis-cli HGET module/samba1/environment SVCPASS
    # output: "Random,1234"
    api-cli run set-backend --agent module/ldapproxy3 --data - <<EOF
    {
        "backend": "samba1",
        "schema": "ad",
        "host": "127.0.0.1",
        "port": 636,
        "tls": true,
        "tls_verify": false,
        "bind_dn": "AD\\ldapservice",
        "bind_password:": "Random,1234"
    }
    EOF
