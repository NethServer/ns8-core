# nextcloud

Start nextcloud with local mariadb and redis instances.
Port `8181` is bound to localhost and accessible using traefik.
Let's Encrypt has been disabled.

Default URL for accessing nextcloud is `https://nextcloud.$(hostname -f)`.

The `build.sh` is the script used to initially generate the systemd units.


## Backup and restore

Assumptions:
- there is a restic rest server is running on 127.0.0.1:8383
- the restic repository already exists
- the restic password is saved inside a `restic_password` file

To execute restic server run as root (the server runs without authentication):
```
podman run -d -p 127.0.0.1:8383:8000 -v restic-server-backup:/data --name rest_server -e OPTIONS="--no-auth --debug" docker.io/restic/rest-server
```

To create the repository execute:
```
echo "nextcloud" > restic_password
podman run --network=host --rm -i -v $HOME/restic_password:/pass -e RESTIC_REPOSITORY=rest:http://127.0.0.1:8383/nextcloud docker.io/restic/restic -p /pass init
```

To execute the backup:
```
./pre-backup && ./backup && ./post-backup
```

To execute the restore:
```
./restore
```
