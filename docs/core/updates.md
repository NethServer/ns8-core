---
layout: default
title: Updates
nav_order: 12
parent: Core
---

# Core updates

The core update changes both the core images and the core modules. In the
first stage cluster nodes update the core components. In the second stage
the leader node launches the update of individual core modules.

## First stage

The first stage applies enhancements and bug fixes for the `core` image.
The core image contains: agents, UI, additional core images and common
actions.

To execute the core update, use:

    api-cli run update-core --data '{"core_url":"ghcr.io/nethserver/core:latest","nodes":[1, 2]}'

The cluster agent running on the leader node forwards the core-update
request to selected nodes.

Each node agent compares the image tag of the `core_url` parameter with
core image version installed on the node. One of the following cases
applies:

- Both tags are valid Semver tags: the image is installed only if
  `core_url` is greater than the installed one.

- One or both tags are invalid Semver tags: `core_url` is pulled only
  if not already present in Podman storage and the update continues

- The `force:true` flag is passed: `core_url` is pulled and update runs
  ignoring the Semver tag check.

When `core_url` is downloaded, additional images listed in the `core`
image label `org.nethserver.images` are downloaded as well.

After all downloads are successful, the `core` image is installed,
replacing the previous one.

Then the node agent runs executable scripts under
`/var/lib/nethserver/node/update-core.d/`. Execution order is
alphabetical. If a script has non-zero exit code, a warning is printed and
execution continues with the next script.

If the update fails on some node, the whole process is aborted by the
leader node at this point.

## Second stage

If every node completes the first stage successfully, the leader node
starts the second stage.

The cluster agent on the leader node runs executable scripts under
`/var/lib/nethserver/cluster/update-core-pre-modules.d/`, as described for
the previous stage. Scripts are obtained from the newly installed core
image.

Then the cluster agent searches installed core modules. For example, it
runs the update-module action for any instance of:

- ldapproxy
- traefik
- samba, openldap (account providers)
- loki

See [module updates]({{site.baseurl}}/modules/updates) for more info about
the module update procedure, that is applied also to core modules.

Finally, the cluster agent runs the scripts under
`/var/lib/nethserver/cluster/update-core-post-modules.d/`, as described in
the previous stage. Scripts are obtained from the newly installed core
image.

The exit code of the `update-core` action is non-zero if any core module
has failed its update.
