# list-domain-users output Schema

```txt
http://schema.nethserver.org/cluster/list-domain-users-output.json
```

List users of a given accounts domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-domain-users-output.json](cluster/list-domain-users-output.json "open original schema") |

## list-domain-users output Type

`object` ([list-domain-users output](list-domain-users-output.md))

## list-domain-users output Examples

```json
{
  "users": [
    {
      "user": "Administrator",
      "display_name": "Administrator"
    }
  ]
}
```

# list-domain-users output Properties

| Property        | Type    | Required | Nullable       | Defined by                                                                                                                                                      |
| :-------------- | :------ | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [users](#users) | `array` | Required | cannot be null | [list-domain-users output](list-domain-users-output-properties-users.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/properties/users") |

## users



`users`

* is required

* Type: `object[]` ([A user descriptor](list-domain-users-output-defs-a-user-descriptor.md))

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-properties-users.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/properties/users")

### users Type

`object[]` ([A user descriptor](list-domain-users-output-defs-a-user-descriptor.md))

# list-domain-users output Definitions

## Definitions group user

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user"}
```

| Property                       | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                      |
| :----------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [user](#user)                  | `string`  | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-user.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/user")                 |
| [display\_name](#display_name) | `string`  | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-display-name.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/display_name") |
| [locked](#locked)              | `boolean` | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-locked.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/locked")             |

### user



`user`

* is required

* Type: `string`

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-user.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/user")

#### user Type

`string`

#### user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### display\_name

Preferred name of a person to be used when displaying entries (RFC 2798/2.3)

`display_name`

* is required

* Type: `string` ([Display name](list-domain-users-output-defs-a-user-descriptor-properties-display-name.md))

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-display-name.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/display_name")

#### display\_name Type

`string` ([Display name](list-domain-users-output-defs-a-user-descriptor-properties-display-name.md))

### locked

True, if the user account has been locked, preventing the user from logging in

`locked`

* is required

* Type: `boolean` ([Locked](list-domain-users-output-defs-a-user-descriptor-properties-locked.md))

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-locked.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/locked")

#### locked Type

`boolean` ([Locked](list-domain-users-output-defs-a-user-descriptor-properties-locked.md))
