# Untitled string in create-cluster input Schema

```txt
http://schema.nethserver.org/cluster/create-cluster-input.json#/properties/endpoint
```

The public WireGuard VPN endpoint in `hostname:port` form. The given port may differ from 55820 as it depends on external configurations (i.e. port-forwards)

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [create-cluster-input.json\*](cluster/create-cluster-input.json "open original schema") |

## endpoint Type

`string`

## endpoint Examples

```json
"my.example.com:55820"
```

```json
"1.2.3.4:60000"
```
