# ns8-scratchpad

NethServer 8 experiments using containers.

Tested on Fedora 33.

Podman running in rootless mode:
- every container has its own user
- traefik reads the dynamic configuration from redis
- containers talk to each other using host network

## Install

Execute as root:
```
dnf install git -y
git clone https://github.com/DavidePrincipi/ns8-scratchpad.git
./setup.sh
```

Access to redis is available from host network:
```
podman run -it --network host  --rm redis redis-cli
```

### Dokuwiki

The script will start a dokuwiki instance with valid SSL certificate, persistence and redirection from HTTP to HTTPs
Default host for the dokuwiki is ``dokuwiki.<fqdn>``, make sure to have a valid DNS public record for it.

Execute as root:
```
./dokuwiki/init.sh
```

