# Netdata

Start and configure a Netdata instance.
The module uses [Netdata official image](https://hub.docker.com/r/netdata/netdata).

## Install and configure

Netdata can be started as:

- collector: collect data from local node and accept data from other nodes
- streamer: collect data from local node and send them to a remote node

The collector node will be the one with a Web UI, netdata will display data from all nodes.

This guide will configure 2 netdata instances:
- the collector on node 1
- the streamer on node 2

### Collector

Instantiate the module, execute on the leader:
```
add-module netdata 1
```

The output of the command will return the instance name.
Output example:
```
Extracting container filesystem ui to /var/lib/nethserver/cluster/ui/apps/netdata19
ui/index.html
9b01c92ae4f28e1ff322803e403b20617c1d40b9fd29c96740d02a87575e549f
{"module_id": "netdata1", "image_name": "Netdata", "image_url": "ghcr.io/nethserver/netdata:latest"}
```

To start a collector instance, launch `configure-module`, by setting the `role` parameter to `collector`:
```
redis-cli LPUSH module/netdata1/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data":{"role":"collector"}}'
```

Finally, setup traefik to access netdata.
Launch `set-host`, by setting the following parameters:
- the module instance name
- the listen URL
- the virtual host name
- the option to enable or disable Let's Encrypt certificate
- the option to enable or disable HTTP to HTTPS redirection

Netdata always listens to port `19999`.
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "set-host", "data": {"instance": "netdata1", "url": "http://127.0.0.1:19999", "host": "netdata.test.local", "lets_encrypt": true, "http2https": true} }'
```

### Streamer

Create a new instance of netdata inside the second node. Execute this on the leader node:
```
add-module netdata 2
```

The ouput will look like this one:
```
Extracting container filesystem ui to /var/lib/nethserver/cluster/ui/apps/netdata2
ui/index.html
bdb872016bb8fa8a81306315ac34e54b25f99aa45f2e3a45286c044a84248d0b
{"module_id": "netdata2", "image_name": "Netdata", "image_url": "ghcr.io/nethserver/netdata:latest"}
```

The `configure-module` task will generate an `api_key` which will be available inside redis as `API_KEY` key.

To start a streamer instance, launch `configure-module`, by setting the following parameter:
- `role`: set to `streamer`
- `api_key`: the collector API key
- `collector_ip`: the collector VPN IP

You can retrieve the collector API key by executing this command:
```
redis-cli HGET module/netdata1/environment API_KEY
```

You can retrieve the collector IP address by executing this command:
```
redis-cli HGET node/1/vpn ip_address
```

Then, execute the configuration for the streamer:
```
redis-cli LPUSH module/netdata2/tasks '{"id": "'$(uuidgen)'", "action": "configure-module", "data":{"role":"streamer", "api_key": "7a9e1687-9fa1-4685-bac3-25c45e6b21e2", "collector_ip": "10.5.4.1"}}'
```

The streamer doesn't require web access: all data are visible through the collector node.

## Uninstall

To uninstall the instance:
```
redis-cli LPUSH module/traefik1/tasks '{"id": "'$(uuidgen)'", "action": "delete-host", "data": {"instance": "netdata1"}}'
remove-module netdata1 --no-preserve
```
