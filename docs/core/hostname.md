---
layout: default
title: hostname FQDN
nav_order: 14
parent: Core
---

# Hostname FQDN

Retrieve the fully qualified domain name of the server by using a python library. In case of exception the function output a generic FQDN `myserver.example.com`

Python example code:

```python
#!/usr/bin/env python3

import sys
import json
import agent

fqdn = agent.get_hostname()
json.dump(fqdn, fp=sys.stdout)
```

output example:
```json
{"hostname": "foo.domain.com"}
```
