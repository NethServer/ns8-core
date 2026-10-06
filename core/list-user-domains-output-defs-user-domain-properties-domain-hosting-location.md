# Domain hosting location Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/location
```

Type of domain hosting. Set to `internal` if the domain is hosted by the cluster, `external` if the domain is provided by a remote service

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-user-domains-output.json\*](cluster/list-user-domains-output.json "open original schema") |

## location Type

`string` ([Domain hosting location](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md))

## location Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value        | Explanation |
| :----------- | :---------- |
| `"internal"` |             |
| `"external"` |             |
