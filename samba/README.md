# samba

The Samba module allows to install one or more AD domains in a NS8
cluster. An AD domain (or realm, in Kerberos terms), can have one or more
domain controllers (DCs). Provisioning a new domain and joining an
additional domain controller to an existing domain are both implemented by
the same action: `provision-domain`.

Do not assign the IP address of a untrusted network! The IP address
assigned to Samba DC exposes internal services that are not designed to be
exposed publicly. The IP address is added to the cluster VPN routes, so
DCs can see each other and perform replication.

The LDAP service of Samba DCs does not allow clear-text LDAP binds: enable TLS or
use Kerberos/GSSAPI authentication. The default node Ldapproxy instance is configured
to enable TLS automatically.

## Provision

Install & provision a new domain:

    # pip install httpie
    add-module samba 1
    export TOKEN=$(http :8080/api/login username=admin password=Nethesis,1234 | jq -r .token)
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

## Configure Ldapproxy

When a Samba module is provisioned for the first time on a node, it checks if the
node default Ldapproxy instance is not bound to another account provider. If so,
Samba invokes the `set-backend` on Ldapproxy automatically.

This command manually binds `ldapproxy1` with the `samba2` local account provider LDAP backend.
TLS is required by Samba (clear-text LDAP binds are not allowed). 

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
    # retrieve the TCP port used by ldapproxy1
    redis-cli HGET module/ldapproxy1/environment TCP_PORT
    # output: "20000"


## Test with ldapsearch

    podman run --network=host --replace --name=alpine -d alpine sh -c 'sleep INF'
    podman exec -ti alpine sh
    # In alpine:
    apk add openldap-clients
    ldapsearch -D 'AD\administrator' -w Nethesis,1234 -H ldap://127.0.0.1:20000 -s sub -b 'DC=ad,DC=dp,DC=nethserver,DC=net' samaccountname=administrator
