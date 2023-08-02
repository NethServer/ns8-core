---
layout: default
title: Validation framework
nav_order: 4
parent: Core
---

# Validation framework

All [agent actions]({{site.baseurl}}/core/agents) inherit a validation system.
Each agent action can include some validation steps using JSON schema or create
custom validation steps.

* TOC
{:toc}

## JSON schema

Basic validation of input and output can be enforced using [JSON
Schema](https://json-schema.org/).

To validate the input just drop a `validate-input.json` file inside the
action directory, to validate the output just create a
`validate-output.json` file.

As general rule, schema files are loaded from special paths:
- `${AGENT_INSTALL_DIR}/validator-definitions.json`: this is always picked
  up by the agent from all actions. It can contain the common data format
  definitions
- `${AGENT_INSTALL_DIR}/actions/<action-name>/validate-input.json`: schema
  applied to validate the input data
- `${AGENT_INSTALL_DIR}/actions/<action-name>/validate-output.json`:
  schema applied to validate the output data

## Custom validation

The JSON schema can perform formal validation of input and output parameters.
Semantic validation, like searching for duplicates or check if a network host is reachable,
should be implemented inside custom scripts.

Additional validation logic can be implemented in an early step script
like `01validation`. If the validation fails, the step must execute the
[set-status](#set-status) command and return a non-zero exit code.

It is recommended to write custom scripts using Python, but any executable script inside
the action directory would do the job.

### Progress handling

Validations steps will block the task progress to 0 inside the UI.
After the validation has passed, the task will be set as running and the progress will be
tracked by the user interface.

When writing Python-based validation scripts, use this snippet:
```python
agent.set_weight(os.path.basename(__file__), 0) # Validation step, no task progress at all
```

### Port validation

we search if the TCP port is used or not by a service return true if the port is used, false is the port is not used. The code below retrieve from the json object of the UI

```python
#!/usr/bin/env python3

#
# Copyright (C) 2022 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import os
import sys
import json
import agent

# retrieve json data
data = json.load(sys.stdin)
port = data['sftp_tcp_port']

# set error validation
agent.set_weight(os.path.basename(__file__), 0)

# if we have the same port, lets go further
if str(port) == os.environ.get('SFTP_TCP_PORT',""):
    sys.exit(0)

if agent.tcp_port_in_use(port):
    agent.set_status('validation-failed')
    json.dump([{'field':'sftp_tcp_port','parameter':'sftp_tcp_port','value':port,'error':'tcp_port_already_used'}],fp=sys.stdout)
    sys.exit(1)
```
