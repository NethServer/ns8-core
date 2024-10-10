---
layout: default
title: Code snippets
nav_order: 140
parent: Modules
---

# Code snippets

Some code examples to ease the life of module developers

## Hostname FQDN

Retrieve the fully qualified domain name of the server by using the `agent` Python library. In case of exception the function outputs the default FQDN `myserver.example.org`.

Python example code:

```python
import agent

fqdn = agent.get_hostname()
print(fqdn)
```

output example:
```
foo.domain.com
```

## Port validation

Check if the TCP port is already used by a service. Exit with an error if the port is used.

```python
import sys
import agent

port=80
if agent.tcp_port_in_use(port):
    sys.exit(1)
```


## Web route validation

Check if HTTP routes are handled by Traefik. The function issues an HTTP request to 127.0.0.1 port 80, setting the `Host` header in the HTTP request and/or the URL path. Any HTTP non-404 response assumes the HTTP route exists.

### Test if a domain-based route is used (domain.com)
```python
import sys
import agent

hostname = 'foo.com'
if agent.http_route_in_use(domain=hostname):
    sys.exit(2)
```

### Test if a path-based route is used (/path)
```python
import sys
import agent
path ='/path' # path fragment of the URL to check
if agent.http_route_in_use(domain=none, path=path):
    sys.exit(2)
```

## Extra python libraries

If the module requires extra Python libraries, you can install them inside the `create-module` action.

Here's an example of a `create-module` action step named `20initialize`:

```bash
#!/bin/bash

pip install bcrypt==4.1.2
```

You can then use the imported `bcrypt` library in all your action steps. A similar script can be added under `update-module.d/`.
