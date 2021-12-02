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
TLS or use Kerberos/GSSAPI authentication in external applications.
Applications that run as cluster modules must connect to the LDAP service
of Samba DCs through the Ldapproxy instance running on the local node; read
the core documentation for more information about Ldapproxy.


## Provision

Provision a new domain:

    api-cli run configure-module --agent module/samba1 --data - <<EOF
    {
        "adminuser":"administrator",
        "adminpass":"Nethesis,1234",
        "realm":"ad.$(hostname -d)",
        "nbdomain":"DP",
        "hostname":"dc1",
        "ipaddress":"10.133.0.2"
    }
    EOF

Further Samba instances for the same `realm` are **joined** to the existing domain.
The command is similar.

    api-cli run configure-module --agent module/samba2 --data - <<EOF
    {
        "adminuser":"administrator",
        "adminpass":"Nethesis,1234",
        "realm":"ad.$(hostname -d)",
        "nbdomain": null,
        "hostname":"dc2",
        "ipaddress":"10.124.0.2"
    }
    EOF

Note that `nbdomain` can be `null` in this case.

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

## Create user

To create a new user, just execute:
```
podman exec -ti samba1 samba-tool user create goofy Nethesis,1234 --given-name=Goofy --surname=Goof --mail=goofy@mail.org
```
