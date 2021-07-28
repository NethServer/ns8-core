# ldapproxy

The Ldapproxy module is a communication broker between an account provider
LDAP backend (like Samba DC, OpenLDAP...) and a consumer module that wants
to authenticate its users against the LDAP service (e.g. Nextcloud).

## Connect to the LDAP service proxy

A default Ldapproxy instance is installed as a core module on each cluster
node. Additional instances can be configured too.

The consumer module runs on the same node where the Ldapproxy instance is
running. To connect with the Ldapproxy default instance:

1. Retrieve the default instance ID

       GET node/${NODE_ID}/default_instance/ldapproxy (=> ldapproxy1)

2. Retrieve the TCP port

       HGET module/ldapproxy1/environment LDAPPORT (=> 20000)

Point the client to `ldap://127.0.0.1:20000`.

If the consumer module LDAP client runs in a Podman container with a
`private` network, add the following arguments to the `podman run`
command:

    --network=slirp4netns:allow_host_loopback=true

Then point the client to `ldap://10.0.2.2:20000` (that IP is routed by
Podman the loopback device on the root network namespace).

## Environment variables

The Ldapproxy environment provides OpenLDAP-compatible variables, like
`LDAPPORT` `LDAPHOST` `LDAPBINDDN` ... Refer also to LDAP.CONF(5) man page
for details.

## Test with ldapsearch

In this test the complete `ldapproxy1` instance environment is exported to
a container where `ldapsearch` runs.

With `host` network:

    podman run --env-file ~ldapproxy1/.config/state/environment --network=host --replace --name=alpine-ldapclient --rm -d alpine sh -c 'apk add openldap-clients ; sleep INF'
    podman exec -i alpine-ldapclient sh -c 'ldapsearch -x -w ${BIND_PASSWORD} -s base'

With `private` network (note the `LDAPURI` override):

    podman run --env-file ~ldapproxy1/.config/state/environment --env LDAPURI=ldap://10.0.2.2:20000 --network=slirp4netns:allow_host_loopback=true --replace --name=alpine-ldapclient --rm -d alpine sh -c 'apk add openldap-clients ; sleep INF'
    podman exec -i alpine-ldapclient sh -c 'ldapsearch -x -w ${BIND_PASSWORD} -s base'

## Set the proxy LDAP backend

Modules with role `accountprovider` can invoke the `set-backend` action on
Ldapproxy instances.

The `set-backend` action configures and restarts the L4 proxy service.

    # pip install httpie
    export TOKEN=$(http :8080/api/login username=admin password=Nethesis,1234 | jq -r .token)
    http :8080/api/module/ldapproxy2/tasks "Authorization: Bearer $TOKEN" <<EOF
    {
        "action":"set-backend", 
        "data": {
            "backend": "samba2",
            "schema": "ad",
            "base_dn": "DC=ad,DC=example,DC=com",
            "host": "127.0.0.1",
            "port": 636,
            "tls": true,
            "tls_verify": false,
            "bind_dn": "ldapservice@ad.example.com",
            "bind_password": "supersecret"
        }
    }
    EOF
