# Leader public VPN UDP endpoint Schema

```txt
http://schema.nethserver.org/cluster/add-node-output.json#/properties/leader_endpoint
```

WireGuard public UDP endpoint address to contact the leader node. Syntax is `hostname:port` or `IP:port`

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-node-output.json\*](cluster/add-node-output.json "open original schema") |

## leader\_endpoint Type

`string` ([Leader public VPN UDP endpoint](add-node-output-properties-leader-public-vpn-udp-endpoint.md))

## leader\_endpoint Examples

```json
"my.example.com:55820"
```

```json
"1.2.3.4:55820"
```
