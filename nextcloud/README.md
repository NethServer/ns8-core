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
{"rootfull": false, "mid": "nextcloud1"}
```

Wait for `add-module` to complete by looking inside `journalctl`.

## Configure

Let's assume that the nextcloud istance is named `nextcloud1`.

Then launch `configure-module`, by setting the following parameters:
- administrator user
- administrator password
- comma-separated list of trused domains

All parameters must be set inside the `data` field separated by a space and terminated with `\n`.

Example:
```
redis-cli LPUSH module/nextcloud1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data": "admin Nethesis,1234 nextcloud.test.loc,localhost\n"}'
```

Finally, setup traefik to access.
Then launch `set-host`, by setting the following parameters:
- the module instance name
- the listen URL
- the virtual host name
- the option to enable (1) or disable (0) Let's Encrypt certificate
- the option to enable (1) or disable (0) HTTP to HTTPS redirection

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
remove-module nextcloud1
```
