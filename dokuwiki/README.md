# Dokuwiki

Start and configure a Dokuwiki instance.
The module uses [Bitnami DokuWiki image](https://github.com/bitnami/bitnami-docker-dokuwiki).

## Install

Instantiate the module:
```
add-module dokuwiki 1
```

The output of the command will return the instance name.
Output example:
```
{"rootfull": false, "mid": "dokuwiki1"}
```

Wait for `add-module` to complete by looking inside `journalctl`.

## Configure

Let's assume that the dokuwiki istance is named `dokuwiki1`.

Then launch `configure-module`, by setting the following parameters:
- wiki name
- administrator user
- administrator password
- administrator mail address

All parameters must be set inside the `data` field separated by a space and terminated with `\n`.

Example:
```
redis-cli LPUSH module/dokuwiki1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data": "MyWiki admin mypassword admin@mydomain.org\n"}'
```

Finally, setup traefik to access.
Then launch `set-host`, by setting the following parameters:
- the module instance name
- the listen URL
- the virtual host name

All parameters must be set inside the `data` field separated by a space and terminated with `\n`.

Example:
```
source /home/dokuwiki1/.config/state/environment
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": "dokuwiki1 http://127.0.0.1:'${TCP_PORT}' mywiki.myhost.org 1 1\n"}'
```

## Uninstall

To uninstall the instance:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": "dokuwiki1\n"}'
remove-module dokuwiki1
```
