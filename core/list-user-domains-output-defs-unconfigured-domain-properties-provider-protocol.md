# Provider protocol Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/protocol
```

Protocol used to communicate with the domain providers.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-user-domains-output.json\*](cluster/list-user-domains-output.json "open original schema") |

## protocol Type

`string` ([Provider protocol](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md))

## protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |
