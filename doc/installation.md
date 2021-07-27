# Installation & First steps

**NOTE**: this is just a prototype, *do not* use in production.

## Core installation

Pick your preferred distribution between Fedora 34 and Debian 11, install it and make sure it's up to date.

Start the installation procedure as `root`:
```
curl https://raw.githubusercontent.com/NethServer/ns8-scratchpad/main/core/install.sh | bash
```

The procedure must be completed by invoking additional commands: follow on-screen instructions.
Run either new cluster initialization (`create-cluster`) or joining an existing cluster (`join-cluster`).

### Install a development branch

Developers may prefer to run `install.sh` with one or more images from a
development branch. The first argument selects the branch name, following
arguments the module names to be pulled from that branch.

    bash install.sh mybranch module1 module2 ...

## First steps

After the install and the first cluster configuration, you can install a new application
or mess around with the core ones.

An experimental UI is also available ad `https://<server_fqdn>/cluster-admin/`:

- default user: `admin`
- default password: `Nethesis,1234`

If all Let's Encrypt requirements for [HTTP-01 challenge](https://letsencrypt.org/docs/challenge-types/#http-01-challenge) are met,
it's possible to request a valid certificate for the FQDN.

Example:
```
pip install httpie
export TOKEN=$(http :8080/api/login username=admin password=Nethesis,1234 | jq -r .token)
http :8080/api/module/traefik5/tasks "Authorization: Bearer $TOKEN" <<EOF
{
    "action":"set-certificate",
    "data": {
        "fqdn": "server.nethserver.org"
    }
}
EOF
```

### Application installation

To install an instance of an available application use:
```
add-module <module> <node_id>
```
The  above command will try to install the latest stable version of the module.

Example to install Dokuwiki on node 1:
```
add-module dokuwiki 1
```

You can also install an image which is not still available inside the repository by using
its URL.

Example to install Dokuwiki directly from the image registry:
```
add-module ghrc.io/nethserver/dokuwiki:mydev 1
```

Many applications need a configuration step after install, for more info, 
please refer to the README of each application.

Available applications:

- [Dokuwiki](../dokuwiki/README.md)
- [Nextcloud](../nextcloud/README.md)
- [Mail](../netdata/README.md)

### Core applications

Core applications installed by default:
- [Redis](#redis)
- [Traefik](../traefik/README.md)

Available core applications:
- [LDAP proxy](../ldapproxy/README.md)
- [Loki](../loki/REDME.md)
- [Promtail](../promtail/README.md)
- [Netdata](../netdata/README.md)
- [Samba](../samba/README.md)


## Uninstall

The `uninstall.sh` script attempts to stop and erase core components and
additional modules. Handle it with care because it erases everything under `/home/*` and `/var/lib/nethserver/*`!

    bash /var/lib/nethserver/node/uninstall.sh

---
Next: [Technical details](details.md)
