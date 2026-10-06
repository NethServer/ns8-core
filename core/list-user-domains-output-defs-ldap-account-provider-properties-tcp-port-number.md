# TCP port number Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/port
```

Can be `null` if the provider is not configured properly

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-user-domains-output.json\*](cluster/list-user-domains-output.json "open original schema") |

## port Type

`integer` ([TCP port number](list-user-domains-output-defs-ldap-account-provider-properties-tcp-port-number.md))

## port Constraints

**minimum**: the value of this number must greater than or equal to: `1`
