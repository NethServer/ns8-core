# nextcloud

Start nextcloud with local mariadb and redis instances.
Port `8181` is bound to localhost and accessible using traefik.
Let's Encrypt has been disabled.

Default URL for accessing nextcloud is `https://nextcloud.$(hostname -f)`.

The `build.sh` is the script used to initially generate the systemd units.
