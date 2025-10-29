# Untitled object in get-domain-group output Schema

```txt
http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-domain-group-output.json\*](cluster/get-domain-group-output.json "open original schema") |

## user Type

`object` ([Details](get-domain-group-output-defs-user.md))

# user Properties

| Property                       | Type     | Required | Nullable       | Defined by                                                                                                                                                                                      |
| :----------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [user](#user)                  | `string` | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-user-properties-user.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/user")                 |
| [display\_name](#display_name) | `string` | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-user-properties-display-name.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/display_name") |

## user



`user`

* is required

* Type: `string`

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-user-properties-user.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/user")

### user Type

`string`

### user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## display\_name

Preferred name of a person to be used when displaying entries (RFC 2798/2.3)

`display_name`

* is required

* Type: `string` ([Display name](get-domain-group-output-defs-user-properties-display-name.md))

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-user-properties-display-name.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/display_name")

### display\_name Type

`string` ([Display name](get-domain-group-output-defs-user-properties-display-name.md))
