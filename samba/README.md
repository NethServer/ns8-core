# samba

## new-domain

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

## Configure ldapproxy

When a Samba module is installed for the first time on a node, it binds itself with the
node default Ldapproxy instance.

This command binds `ldapproxy1` with the new `samba1` local account provider LDAP backend.
TLS is required by Samba (clear-text LDAP binds are not allowed). 

    http :8080/api/module/ldapproxy1/tasks "Authorization: Bearer $TOKEN" <<EOF
    {
        "action":"set-backend", 
        "data": {
            "backend": "samba1",
            "schema": "ad",
            "host": "127.0.0.1",
            "port": 636,
            "tls": true
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

## join-domain

Join an additional DC (installed on another node):

1. Complete the provisioning of the first instance

2. Proceed with:

       http :8080/api/module/samba2/tasks "Authorization: Bearer $TOKEN" <<EOF
       {
           "action":"provision-domain",
           "data":{
               "adminuser":"administrator",
               "adminpass":"Nethesis,1234",
               "ipaddress":"10.134.0.56"
           }
       }
       EOF
