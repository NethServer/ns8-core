# rsync

This directory contains scripts for the **rsync** image. It provides both
client and server configuration to transfer a module instance state to
another instance.

Invocation (server)

    podman run --rm --privileged --network=host --workdir=/srv --env=RSYNCD_NETWORK=10.5.4.0/24 --env=RSYNCD_ADDRESS=cluster-localnode --env=RSYNCD_PORT=20024 --env=RSYNCD_USER=dokuwiki1 --env=RSYNCD_PASSWORD=534421aa2c25f-d383-498f-8d77-f33c41fc01fb --volume=/dev/log:/dev/log --volume=/home/dokuwiki7/.config/state:/srv/state --volume=dokuwiki-data:/srv/volumes/dokuwiki-data ghcr.io/nethserver/rsync:latest

Invocation (client)

    podman run -i --workdir=/srv --rm --network=host --privileged --env=RSYNC_PASSWORD=534421aa2c25f-d383-498f-8d77-f33c41fc01fb --volume=/home/dokuwiki1/.config/state:/srv/state --volume=dokuwiki-data:/srv/volumes/dokuwiki-data ghcr.io/nethserver/rsync:latest rsync -a --info=progress2 --delete-after --filter=". -" ./ rsync://dokuwiki1@10.5.4.1:20024/state/

## Basic Podman flags

- `--network=host`: the rsyncd server must be accessible from other nodes
  of the cluster, see also the RSYNCD_ADDRESS environment variable.
- `--privileged`: both the client and the server need full access to volume contents

## Volume mounts

When the server is started, mount the following volumes:

- `--volume=/dev/log:/dev/log`: to send log messages to syslog/journald with RSYNCD_SYSLOG_TAG
- `--volume=$AGENT_STATE_DIR:/srv/state`: `state/` directory is always mounted, like backups
- `--volume=volname:/srv/volumes/volname`: mount the volumes under `/srv/volumes`, like backups

For the client side, /dev/log is not necessary.

## Environment variables

Server

- RSYNCD_ADDRESS: rsyncd listen address
- RSYNCD_PORT: rsyncd listen port
- RSYNCD_USER: user name for rsync authentication
- RSYNCD_PASSWORD: password for rsync authentication
- RSYNCD_NETWORK: allow connection only if client comes from the given network
- RSYNCD_SYSLOG_TAG: set the syslog identity

Client

- RSYNC_PASSWORD: client password, username is specified in the command arguments