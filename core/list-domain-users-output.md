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
      "display_name": "Administrator",
      "locked": false,
      "expired": false,
      "mail": "",
      "must_change": false,
      "password_expiration": 1693040546
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

| Property                                     | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                         |
| :------------------------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [user](#user)                                | `string`  | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-user.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/user")                                    |
| [display\_name](#display_name)               | `string`  | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-display-name.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/display_name")                    |
| [locked](#locked)                            | `boolean` | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-locked.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/locked")                                |
| [expired](#expired)                          | `boolean` | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-expired.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/expired")                              |
| [mail](#mail)                                | `string`  | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-email-address.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/mail")                           |
| [must\_change](#must_change)                 | `boolean` | Optional | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-must-change-password.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/must_change")             |
| [password\_expiration](#password_expiration) | `integer` | Required | cannot be null | [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-password-expiration-time.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/password_expiration") |

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

### expired

True, if the user account has expired, preventing the user from logging in

`expired`

* is required

* Type: `boolean` ([Expired](list-domain-users-output-defs-a-user-descriptor-properties-expired.md))

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-expired.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/expired")

#### expired Type

`boolean` ([Expired](list-domain-users-output-defs-a-user-descriptor-properties-expired.md))

### mail

Email address of the user

`mail`

* is required

* Type: `string` ([Email address](list-domain-users-output-defs-a-user-descriptor-properties-email-address.md))

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-email-address.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/mail")

#### mail Type

`string` ([Email address](list-domain-users-output-defs-a-user-descriptor-properties-email-address.md))

### must\_change

True, if the user must change the password at next logon

`must_change`

* is optional

* Type: `boolean` ([Must change password](list-domain-users-output-defs-a-user-descriptor-properties-must-change-password.md))

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-must-change-password.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/must_change")

#### must\_change Type

`boolean` ([Must change password](list-domain-users-output-defs-a-user-descriptor-properties-must-change-password.md))

### password\_expiration

UNIX timestamp when the password will expire, or -1 if it never expires

`password_expiration`

* is required

* Type: `integer` ([Password expiration time](list-domain-users-output-defs-a-user-descriptor-properties-password-expiration-time.md))

* cannot be null

* defined in: [list-domain-users output](list-domain-users-output-defs-a-user-descriptor-properties-password-expiration-time.md "http://schema.nethserver.org/cluster/list-domain-users-output.json#/$defs/user/properties/password_expiration")

#### password\_expiration Type

`integer` ([Password expiration time](list-domain-users-output-defs-a-user-descriptor-properties-password-expiration-time.md))
