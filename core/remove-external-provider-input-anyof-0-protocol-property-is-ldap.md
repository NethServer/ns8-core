# Protocol property is ldap Schema

```txt
http://schema.nethserver.org/cluster/remove-external-provider-input.json#/anyOf/0/not
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-external-provider-input.json\*](cluster/remove-external-provider-input.json "open original schema") |

## not Type

`object` ([Protocol property is ldap](remove-external-provider-input-anyof-0-protocol-property-is-ldap.md))

# not Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                            |
| :-------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [protocol](#protocol) | `string` | Optional | cannot be null | [remove-external-provider input](remove-external-provider-input-anyof-0-protocol-property-is-ldap-properties-protocol.md "http://schema.nethserver.org/cluster/remove-external-provider-input.json#/anyOf/0/not/properties/protocol") |

## protocol



`protocol`

* is optional

* Type: `string`

* cannot be null

* defined in: [remove-external-provider input](remove-external-provider-input-anyof-0-protocol-property-is-ldap-properties-protocol.md "http://schema.nethserver.org/cluster/remove-external-provider-input.json#/anyOf/0/not/properties/protocol")

### protocol Type

`string`

### protocol Constraints

**constant**: the value of this property must be equal to:

```json
"ldap"
```
