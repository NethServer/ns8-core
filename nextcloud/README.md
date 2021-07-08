# Nextcloud

Start and configure a Nexctloud instance:
- with PHP FPM + nginx as a proxy
- redis caching
- MariaDB database
The module uses [Official Nextcloud image](https://hub.docker.com/_/nextcloud).

## Install

Instantiate the module:
```
add-module nextcloud 1
```

The output of the command will return the instance name.
Output example:
```
{"module_id": "nextcloud1", "image_name": "Dokuwiki", "image_url": "ghcr.io/nethserver/nextcloud:latest"}
```

Wait for `add-module` to complete by looking inside `journalctl`.

## Configure

Let's assume that the nextcloud istance is named `nextcloud1`.

Then launch `configure-module`, by setting the following parameters:
- administrator user
- administrator password
- comma-separated list of trused domains

Example:
```
LPUSH module/nextcloud1/tasks  '{"id": "'$(uuidgen)'", "action": "configure-module", "data": {"username": "admin", "password": "admin", "trusted_domains": ["nextcloud.example.org"]}}'
```

Finally, setup traefik to access.
Then launch `set-route`, by setting the following parameters:
- the module instance name
- the listen URL
- the virtual host name
- the option to enable or disable Let's Encrypt certificate
- the option to enable or disable HTTP to HTTPS redirection

See `/home/nextcloid1/.config/state/environment` to find the `TCP_PORT` where the instance will listen. In this exmple `TCP_PORT` is `2000`:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-route", "data": {"instance": "nextcloud", "url": "http://127.0.0.1:20000", "host": "nextcloud.example.org, "lets_encrypt": true, "http2https": true} }'
```

The `TCP_PORT` parameter can be found inside `/home/nextcloud1/.config/state/environment`.

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-route", "data": {"instance": "nextcloud1", "url": "http://127.0.0.1:"'${TCP_PORT}'", "host": "mywiki.myhost.org", "lets_encrypt": true, "http2https": true} }'
```

To execute `occ` command inside an instance:
```
ssh nextcloud1@localhost
podman exec -ti --user www-data nextcloud-app php occ
```



## Uninstall

To uninstall the instance:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": "nextcloud1\n"}'
remove-module nextcloud1 --no-preserve
```
