# nsdc

This module runs as singleton rootfull instance.

- *Rootfull* because Samba needs special privileges to store ACLs in the filesystem extended attributes
- *Singleton* because Samba services are bound to a host IP address, to serve LAN clients

Initialize the Redis DB with

    HSET nsdc0/module.env EVENTS_IMAGE ghcr.io/nethserver/cplane-nsdc:latest NSDC_IMAGE ghcr.io/nethserver/nsdc:latest IPADDRESS 10.133.0.5 HOSTNAME nsdc1.ad.dp.nethserver.net NBDOMAIN AD REALM AD.DP.NETHSERVER.NET ADMINPASS Nethesis,1234

Then start the module installation:
```
podman run -i --network host --rm docker.io/redis:6-alpine redis-cli <<EOF
PUBLISH $(hostname -s):module-rootfull.init nsdc0
EOF
```
