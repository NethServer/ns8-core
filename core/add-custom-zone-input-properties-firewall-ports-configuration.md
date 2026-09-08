# Firewall ports configuration Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/ports
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-custom-zone-input.json\*](node/add-custom-zone-input.json "open original schema") |

## ports Type

`array` ([Firewall ports configuration](add-custom-zone-input-properties-firewall-ports-configuration.md))

## ports Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

## ports Examples

```json
"25/tcp"
```

```json
"21/udp"
```

```json
"9000-9010/tcp"
```
