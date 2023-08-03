---
layout: default
title: Code snippets
nav_order: 14
parent: Modules
---

# Code snippets

Some code examples to ease the life of module developers

## Hostname FQDN

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

## Port validation

we search if the TCP port is used or not by a service return true if the port is used, false is the port is not used.

```python
#!/usr/bin/env python3

#
# Copyright (C) 2022 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import sys
import agent

port=80
if agent.tcp_port_in_use(port):
    sys.exit(1)
```


## Web route validation

we search if the domain (foo.com) or the webpath (/foo) is used or not by a service return true if the route is used, false is the route is not used.

### test if the domain is used (domain.com)
```python
#!/usr/bin/env python3

#
# Copyright (C) 2022 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import sys
import agent

# Setup default values
hostname = 'foo.com'

if agent.http_route_in_use(domain=hostname):
    sys.exit(2)
```

### test if the web path is not used (/path)
```python
#!/usr/bin/env python3

#
# Copyright (C) 2022 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import sys
import agent

# Setup default values
path ='/path'

if agent.http_route_in_use(domain=none, path=path):
    sys.exit(2)
```
