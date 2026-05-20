# Ports protocol Schema

```txt
http://schema.nethserver.org/node/allocate-ports-input.json#/properties/protocol
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                           |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [allocate-ports-input.json\*](node/allocate-ports-input.json "open original schema") |

## protocol Type

`string` ([Ports protocol](allocate-ports-input-properties-ports-protocol.md))

## protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value   | Explanation |
| :------ | :---------- |
| `"tcp"` |             |
| `"udp"` |             |
