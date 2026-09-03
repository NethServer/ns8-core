# Endpoint UDP port for VPN Schema

```txt
http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_port
```

The UDP port number where Wireguard is listening. It must be reachable by other cluster nodes

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [promote-node-input.json\*](cluster/promote-node-input.json "open original schema") |

## endpoint\_port Type

`integer` ([Endpoint UDP port for VPN](promote-node-input-properties-endpoint-udp-port-for-vpn.md))

## endpoint\_port Constraints

**minimum**: the value of this number must greater than or equal to: `1`
