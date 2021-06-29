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
Trying to pull ghcr.io/nethserver/dokuwiki:latest...
Getting image source signatures
Copying blob sha256:d08608c7bf4e48259ce20524ad0e4285e4e9576a223890ab6092c55d9a03c081
Copying config sha256:9b4a5d22259145875b18c00daf9bf344f57d5acb39b46dd5fa4e89d73e26d6d9
Writing manifest to image destination
Storing signatures
9b4a5d22259145875b18c00daf9bf344f57d5acb39b46dd5fa4e89d73e26d6d9
Extracting container filesystem ui to /var/lib/nethserver/cluster/ui/apps/dokuwiki1
ui/index.html
4234e1db74f05d24beb04feab8038d4dfb1af123abe29a196e1d390b72cb68e0
{"module_id": "dokuwiki1", "image_name": "Dokuwiki", "image_url": "ghcr.io/nethserver/dokuwiki:latest"}
```

## Configure

Let's assume that the dokuwiki istance is named `dokuwiki1`.

Then launch `configure-module`, by setting the following parameters:
- wiki name
- administrator user
- administrator password
- administrator mail address

Example:
```
redis-cli LPUSH module/dokuwiki1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data": {"wiki_name": "mywiki", "username": "admin", "password": "admin", "email": "admin@test.local", "user_full_name": "Administrator"}}'
```

Finally, setup traefik to access.
Then launch `set-route`, by setting the following parameters:
- the module instance name
- the listen URL
- the virtual host name
- the option to enable or disable Let's Encrypt certificate
- the option to enable or disable HTTP to HTTPS redirection

See `/home/dokuwiki1/.config/state/environment` to find the `TCP_PORT` where the instance will listen. In this exmple `TCP_PORT` is `2000`:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-route", "data": {"instance": "dokuwiki1", "url": "http://127.0.0.1:20000", "host": "mywiki.example.org", "lets_encrypt": true, "http2https": true} }'
```

## Uninstall

To uninstall the instance:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": {"instance": "dokuwiki1"}}'
remove-module dokuwiki1 --no-preserve
```
