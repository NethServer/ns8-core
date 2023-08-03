---
layout: default
title: Code snippets
nav_order: 14
parent: Modules
---

# Hostname FQDN

Retrieve the fully qualified domain name of the server by using a python library. In case of exception the function output a generic FQDN `myserver.example.com`

Python example code:

```python
#!/usr/bin/env python3

import agent

fqdn = agent.get_hostname()
print(fqdn)
```

output example:
```
foo.domain.com
```
