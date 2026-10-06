# Transport protocol Schema

```txt
http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/transport
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-service-providers-input.json\*](agent/list-service-providers-input.json "open original schema") |

## transport Type

`string` ([Transport protocol](list-service-providers-input-properties-transport-protocol.md))

## transport Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## transport Examples

```json
"tcp"
```

```json
"udp"
```

```json
"http"
```
