---
layout: default
title: smtp smarthost
nav_order: 14
parent: Core
---

# Smtp smarthost

Settings and credentials for a SMTP smarthost provider are stored inside redis, 
a module can discover it by using a python library. The cluster won't send the 
email for you, the module will do it if it has been designed to send emails. 

Python example code:

```python
#!/usr/bin/env python3
import json
import agent
from agent.smarthost import SmartHost


sp = SmartHost()
smtp = sp.load_smtp()
json.dump(smtp, fp=sys.stdout)
```

output example:
```json
{"port": 587, "host": "", "username": "", "password": "", "enabled": false, "tls": true, "tls_verify": true}
```

- tls: enable startTls (true/false)
- tls_verify: verify if the certificate is valid and if the hostname is associated to the certificate (true/false)


The module can handle account provider changes by defining an event
handler. Create an executable script with path
`${AGENT_INSTALL_DIR}/events/smarthost-provider-changed/10handler` and run
any command from it. For instance:

```shell
mycommand && systemctl --user reload mymodule.service
```
