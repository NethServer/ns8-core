# Dokuwiki

How to start a Dokuwiki instance.

Install the module:
```
add-module dokuwiki 1
```

Configure the module, set the wiki name inside the `data` field (remember to add the `\n` just after the name):
```
redis-cli LPUSH module/dokuwiki1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data": "MyWiki\n"}'
```

Verify the assigned `TCP_PORT` from `/home/dokuwikiX/.config/state/environment`, where `X` is the number of the instance returned from `add-module`.
Use the given port to create a traefik route:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": "dokuwiki1 http://127.0.0.1:$TCP_PORT mywiki.myhost.org\n"}'
```
