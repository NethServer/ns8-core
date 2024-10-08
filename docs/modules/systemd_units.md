---
layout: default
title: Systemd units
nav_order: 60
parent: Modules
---

# Systemd units

A module should have at least one Systemd unit to keep the service running.
The Systemd unit should:
- enable the service on boot
- start and stop podman processes
- ensure restart at least on failure

The Systemd unit can:
- change the working directory to ease finding additional files: `WorkingDirectory=%S/state`
- execute ancillary commands in `Exec*` steps with the `runagent` helper: `ExecStartPre=/usr/local/bin/runagent expand-templates`
- import and use additional environment variables: `EnvironmentFile=%S/state/environment`

When `EnvironmentFile=%S/state/environment` is used in Systemd service
units, remember that the file syntax is not designed to be compatible with
Systemd. Values are stored as raw values, i.e. special characters are not
escaped or quoted:

1. Limit the use of this technique to variables with controlled values, e.g. `IMAGE_URL`.
2. Avoid storing and reading arbitrary strings, like randomly generated
   tokens, and values coming from the user's input.

See [Dokuwiki
unit](https://github.com/NethServer/ns8-dokuwiki/blob/main/imageroot/systemd/user/dokuwiki.service)
as an example.


