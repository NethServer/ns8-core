# ns8-scratchpad

NethServer 8 experiments using containers.

Tested on Fedora 33 and Debian 10.

Podman running in rootless mode:
- every container has its own user
- traefik reads the dynamic configuration from redis
- containers talk to each other using host network

## Install

Execute as root:
```
curl -L -o ns8-scratchpad.tar.gz https://github.com/DavidePrincipi/ns8-scratchpad/archive/main.tar.gz
tar xvzf ns8-scratchpad.tar.gz
cd ns8-scratchpad-main
```

On Fedora:
```
./setup.sh
```

On Debian (manual steps are required at the end):
```
./setup-debian.sh
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

