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
Trying to pull ghcr.io/nethserver/nextcloud:latest...
Getting image source signatures
Copying blob sha256:d08608c7bf4e48259ce20524ad0e4285e4e9576a223890ab6092c55d9a03c081
Copying config sha256:9b4a5d22259145875b18c00daf9bf344f57d5acb39b46dd5fa4e89d73e26d6d9
Writing manifest to image destination
Storing signatures
9b4a5d22259145875b18c00daf9bf344f57d5acb39b46dd5fa4e89d73e26d6d9
Extracting container filesystem ui to /var/lib/nethserver/cluster/ui/apps/nextcloud1
ui/index.html
4234e1db74f05d24beb04feab8038d4dfb1af123abe29a196e1d390b72cb68e0
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
Then launch `set-host`, by setting the following parameters:
- the module instance name
- the listen URL
- the virtual host name
- the option to enable or disable Let's Encrypt certificate
- the option to enable or disable HTTP to HTTPS redirection

See `/home/nextcloid1/.config/state/environment` to find the `TCP_PORT` where the instance will listen. In this exmple `TCP_PORT` is `2000`:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": {"instance": "nextcloud", "url": "http://127.0.0.1:20000", "host": "nextcloud.example.org, "lets_encrypt": true, "http2https": true} }' 
```


All parameters must be set inside the `data` field separated by a space and terminated with `\n`.
The `TCP_PORT` parameter can be found inside `/home/nextcloud1/.config/state/environment`.

Example:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": "nextcloud1 http://127.0.0.1:'${TCP_PORT}' mywiki.myhost.org 1 1\n"}'
```

## Uninstall

To uninstall the instance:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": "nextcloud1\n"}'
remove-module nextcloud1 --no-preserve
```
