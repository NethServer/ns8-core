# Allowed destinations Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1/properties/destinations
```

IPv4 address or network, with an optional port and protocol

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-custom-zone-input.json\*](node/add-custom-zone-input.json "open original schema") |

## destinations Type

`string[]`

## destinations Examples

```json
"10.5.5.0/24:443/tcp"
```

```json
"192.168.1.10:445/tcp"
```

```json
"192.168.1.0/24"
```
