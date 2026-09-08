# Transport protocol Schema

```txt
http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/transport
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                              |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-service-providers-output.json\*](agent/list-service-providers-output.json "open original schema") |

## transport Type

`string` ([Transport protocol](list-service-providers-output-service-discovery-information-properties-transport-protocol.md))

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
