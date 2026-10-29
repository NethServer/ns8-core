# LDAP database schema Schema

```txt
http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/schema
```

The LDAP schema is probed automatically if the value is `null` or the property is missing

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-external-domain-input.json\*](cluster/add-external-domain-input.json "open original schema") |

## schema Type

`string` ([LDAP database schema](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md))

## schema Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value       | Explanation |
| :---------- | :---------- |
| `"ad"`      |             |
| `"rfc2307"` |             |
