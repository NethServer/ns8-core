---
layout: default
title: SMTP smarthost
nav_order: 14
parent: Core
---

# SMTP smarthost

Settings and credentials for a SMTP smarthost provider are stored inside Redis, 
a module can discover it by using a Python library. The cluster won't send the 
email for you, the module will do it if it has been designed to send emails. 

Python example code:

```python
#!/usr/bin/env python3

import sys
import json
import agent


smtp = agent.get_smarthost_settings(agent.redis_connect())
if smtp['enabled'] == True:
    pass
json.dump(smtp, fp=sys.stdout)
```

output example:
```json
{"port": 587, "host": "", "username": "", "password": "", "enabled": false, "encrypt_smtp": "starttls", "tls_verify": true}
```

- encrypt_smtp: select the SMTP encryption method (`none`, `starttls`, `tls`)
- tls_verify: verify if the certificate is valid and if the hostname is associated to the certificate (true/false)
- username: it might be empty if the smtp server does not require an user authentication

The module can handle smarthost setup changes by defining an event
handler. Create an executable script with path
`${AGENT_INSTALL_DIR}/events/smarthost-changed/10handler` and run
any command from it. For instance:

```shell
mycommand && systemctl --user reload mymodule.service
```
