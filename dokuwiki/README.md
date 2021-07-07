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
{"module_id": "dokuwiki1", "image_name": "Dokuwiki", "image_url": "ghcr.io/nethserver/dokuwiki:latest"}
```

## Configure

Let's assume that the dokuwiki istance is named `dokuwiki1`.

Then launch `configure-module`, by setting the following parameters:
- `wiki_name`: wiki name
- `username`: administrator user
- `password`: administrator password
- `email`: administrator mail address
- `user_full_name`: administrator full name
- `host`: a fully qualified domain name for the wiki
- `http2https`: enable or disable HTTP to HTTPS redirection
- `lets_encrypt`: enable or disable Let's Encrypt certificate

Example:
```
redis-cli LPUSH module/dokuwiki1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data":{"wiki_name":"mywiki","username":"admin","password":"admin","email":"admin@test.local","user_full_name":"Administrator","host":"dokuwiki.test.local","http2https":true,"lets_encrypt":false}}'
```

The above command will:
- start and configure the Dokuwiki instance
- configure a virtual host for trafik to access the instance

## Uninstall

To uninstall the instance and remove traefik virtual host:
```
remove-module dokuwiki1 --no-preserve
```
