# Backup & restore

**THIS DOCUMENTATION IS OBSOLETE**

## First configuration

The core and each instance can implement the `backup` event. The event takes the name of the backup as first argument.

First, generate a password for the new backup and set restic repository base URL:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SETNX backup/<backup_name>/password <password>
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SET backup/<backup_name>/base_repository rest:http://127.0.0.1:8383
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SETNX backup/backup1/password $(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli SET backup/backup1/base_repository rest:http://127.0.0.1:8383
```

## Core

**Backup**

To backup the core execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH <hostname>:backup <backup_name>
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH $(hostname -s):backup backup1
```

**Restore**

To restore the core execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH $(hostname -s):restore backup1
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH $(hostname -s):restore backup1
```

## Module instance

**Backup**

To execute the backup of an instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH <user_instance>:backup <backup_name>
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH nextcloud0:backup backup1
```

**Restore**

Each instance implementing the `backup` event, must implement also the `restore` event. The event takes the name of the backup as first argument.

To execute the restore of an instance execute:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH <user_instance>:restore <backup_name>
```

Example:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli PUBLISH nextcloud0:restore backup1
```

