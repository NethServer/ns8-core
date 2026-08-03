# Leader public VPN UDP endpoint Schema

```txt
http://schema.nethserver.org/cluster/join-node-input.json#/properties/leader_endpoint
```

WireGuard public UDP endpoint address to contact the leader node. Syntax is `hostname:port` or `IP:port`

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [join-node-input.json\*](cluster/join-node-input.json "open original schema") |

## leader\_endpoint Type

`string` ([Leader public VPN UDP endpoint](join-node-input-properties-leader-public-vpn-udp-endpoint.md))

## leader\_endpoint Examples

```json
"my.example.com:55820"
```

```json
"1.2.3.4:55820"
```
