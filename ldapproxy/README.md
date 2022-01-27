# ldapproxy

The Ldapproxy module is a communication broker between an account provider
LDAP backend (like Samba DC, OpenLDAP...) and a consumer module that wants
to authenticate its users against the LDAP service (e.g. Nextcloud).

## Connect to the LDAP service proxy

A default Ldapproxy instance is installed as a core module on each cluster
node.

The consumer module runs on the same node where the Ldapproxy instance is
running. See how to connect with Ldapproxy with the [Python
`agent.ldapproxy`
package](https://github.com/NethServer/ns8-scratchpad/blob/main/doc/details.md#users-and-groups-ldapproxy).


## Test with ldapsearch

Annotate the TCP port and bind credentials from the output of `runagent
python3 -magent.ldapproxy`:

    export TCP_PORT=20000 BIND_DN=ldapservice@ad.dp.nethserver.net BIND_PASSWORD=xxx

Test with Podman `host` network:

    podman run --network=host --replace --name=alpine-ldapclient --rm -d alpine sh -c 'apk add openldap-clients ; sleep INF'
    podman exec -i alpine-ldapclient ldapsearch -x -D "${BIND_DN}" -w "${BIND_PASSWORD}" -s base -H ldap://"127.0.0.1:${TCP_PORT}"

Test with Podman `private` network:

    podman run --network=slirp4netns:allow_host_loopback=true --replace --name=alpine-ldapclient --rm -d alpine sh -c 'apk add openldap-clients ; sleep INF'
    podman exec -i alpine-ldapclient ldapsearch -x -D "${BIND_DN}" -w "${BIND_PASSWORD}" -s base -H ldap://"10.0.2.2:${TCP_PORT}"
