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

## Atomic file write

Write a file so that readers never see partial contents, even if the
process crashes or the system loses power during the write. The
`agent.safe_open()` context manager writes to a temporary file in the
same directory as the target, flushes it to disk, then atomically
replaces the target and syncs the directory.

```python
import agent

with agent.safe_open("myapp.conf", "w") as fo:
    fo.write("key=value\n")
```

- Only write modes, like `w` or `wb`, are supported; read and append
  modes raise `ValueError`. The `encoding` argument applies to text
  modes.
- If the target file exists, its permission bits are copied to the new
  file. Otherwise the new file is created with mode `0666` minus the
  process umask.
- If the `with` block raises an exception, the target file is left
  untouched and the temporary file is removed.

The agent library itself uses `safe_open()` to write the `environment`
file in `agent.set_env()` and `agent.write_envfile()`. The core uses it
also to generate the HAProxy configuration of the rclone backup gateway
(`rclonegwctl`).

`safe_open()` does not serialize concurrent writers: the last one wins.
For read-modify-write cycles, wrap the whole cycle with the
`agent.exclusive_file_lock()` context manager, as `agent.set_env()`
does:

```python
import agent

with agent.exclusive_file_lock("myapp.conf"):
    with open("myapp.conf") as fi:
        data = fi.read()
    with agent.safe_open("myapp.conf", "w") as fo:
        fo.write(data + "other=value\n")
```

The lock is an advisory `flock()` on a hidden `.<filename>.lock` file in
the same directory, so it is effective only among processes that use
the same function.

## Extra python libraries

If the module requires extra Python libraries, you can install them inside the `create-module` action.

Here's an example of a `create-module` action step named `20initialize`:

```bash
#!/bin/bash

pip install bcrypt==4.1.2
```

You can then use the imported `bcrypt` library in all your action steps. A similar script can be added under `update-module.d/`.
