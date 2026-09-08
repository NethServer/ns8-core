# Untitled object in get-domain-user output Schema

```txt
http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-domain-user-output.json\*](cluster/get-domain-user-output.json "open original schema") |

## user Type

`object` ([Details](get-domain-user-output-defs-user.md))

# user Properties

| Property                       | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                   |
| :----------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [user](#user)                  | `string`  | Required | cannot be null | [get-domain-user output](get-domain-user-output-defs-user-properties-user.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/user")                                 |
| [display\_name](#display_name) | `string`  | Required | cannot be null | [get-domain-user output](get-domain-user-output-defs-user-properties-display-name.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/display_name")                 |
| [locked](#locked)              | `boolean` | Required | cannot be null | [get-domain-user output](get-domain-user-output-defs-user-properties-locked.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/locked")                             |
| [groups](#groups)              | `array`   | Required | cannot be null | [get-domain-user output](get-domain-user-output-defs-user-properties-list-of-groups-the-user-belongs-to.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/groups") |

## user



`user`

* is required

* Type: `string`

* cannot be null

* defined in: [get-domain-user output](get-domain-user-output-defs-user-properties-user.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/user")

### user Type

`string`

### user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## display\_name

Preferred name of a person to be used when displaying entries (RFC 2798/2.3)

`display_name`

* is required

* Type: `string` ([Display name](get-domain-user-output-defs-user-properties-display-name.md))

* cannot be null

* defined in: [get-domain-user output](get-domain-user-output-defs-user-properties-display-name.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/display_name")

### display\_name Type

`string` ([Display name](get-domain-user-output-defs-user-properties-display-name.md))

## locked

True, if the user account has been locked, preventing the user from logging in

`locked`

* is required

* Type: `boolean` ([Locked](get-domain-user-output-defs-user-properties-locked.md))

* cannot be null

* defined in: [get-domain-user output](get-domain-user-output-defs-user-properties-locked.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/locked")

### locked Type

`boolean` ([Locked](get-domain-user-output-defs-user-properties-locked.md))

## groups



`groups`

* is required

* Type: `object[]` ([Group details](get-domain-user-output-defs-group-details.md))

* cannot be null

* defined in: [get-domain-user output](get-domain-user-output-defs-user-properties-list-of-groups-the-user-belongs-to.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/user/properties/groups")

### groups Type

`object[]` ([Group details](get-domain-user-output-defs-group-details.md))
