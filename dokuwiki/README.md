# Dokuwiki

How to start a Dokuwiki instance.

Install the module:
```
add-module dokuwiki 1
```

Let's assume that the dokuwiki istance is named `dokuwiki1`.

Wait for `add-module` to complete by looking inside `journalctl`, then onfigure the module, 
set the wiki name inside the `data` field (remember to add the `\n` just after the name):
```
redis-cli LPUSH module/dokuwiki1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data": "MyWiki\n"}'
```

Setup traefik routes:
```
source /home/dokuwiki1/.config/state/environment && echo '{"id": "'$(uuidgen)'", "action": "set-host", "data": "docuwiki1 http://127.0.0.1:'${TCP_PORT}' mywiki.myhost.org\n"}'
```
