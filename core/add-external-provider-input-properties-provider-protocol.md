# Provider protocol Schema

```txt
http://schema.nethserver.org/cluster/add-external-provider-input.json#/properties/protocol
```

Protocol used to communicate with the domain providers.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-external-provider-input.json\*](cluster/add-external-provider-input.json "open original schema") |

## protocol Type

`string` ([Provider protocol](add-external-provider-input-properties-provider-protocol.md))

## protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |
