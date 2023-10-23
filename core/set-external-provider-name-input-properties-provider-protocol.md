# Provider protocol Schema

```txt
http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/properties/protocol
```

Protocol used to communicate with the domain providers.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [set-external-provider-name-input.json\*](cluster/set-external-provider-name-input.json "open original schema") |

## protocol Type

`string` ([Provider protocol](set-external-provider-name-input-properties-provider-protocol.md))

## protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |
