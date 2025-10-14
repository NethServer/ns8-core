# alter-user input Schema

```txt
http://schema.nethserver.org/cluster/alter-user-input.json
```

Alter an user account in the Redis DB for the cluster administration web interface

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-user-input.json](cluster/alter-user-input.json "open original schema") |

## alter-user input Type

`object` ([alter-user input](alter-user-input.md))

## alter-user input Examples

```json
{
  "user": "admin",
  "set": {
    "display_name": "The Administrator",
    "password": "New Password"
  },
  "grant": [],
  "revoke": []
}
```

# alter-user input Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                |
| :---------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [user](#user)     | `string` | Required | cannot be null | [alter-user input](cluster-definitions-strict-username-format.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/user")           |
| [set](#set)       | `object` | Required | cannot be null | [alter-user input](cluster-definitions-user-attributes.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/set")                   |
| [grant](#grant)   | `array`  | Optional | cannot be null | [alter-user input](alter-user-input-properties-grant-assertions-list.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/grant")   |
| [revoke](#revoke) | `array`  | Optional | cannot be null | [alter-user input](alter-user-input-properties-revoke-assertions-list.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/revoke") |

## user

A username string can be written in many ways. This is a quite strict form. See <https://www.unix.com/man-page/linux/8/useradd/>

`user`

* is required

* Type: `string` ([Strict username format](cluster-definitions-strict-username-format.md))

* cannot be null

* defined in: [alter-user input](cluster-definitions-strict-username-format.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/user")

### user Type

`string` ([Strict username format](cluster-definitions-strict-username-format.md))

### user Constraints

**maximum length**: the maximum number of characters for this string is: `32`

**minimum length**: the minimum number of characters for this string is: `1`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[a-z_][a-z0-9_-]*[$]?$
```

[try pattern](https://regexr.com/?expression=%5E%5Ba-z_%5D%5Ba-z0-9_-%5D*%5B%24%5D%3F%24 "try regular expression with regexr.com")

### user Examples

```json
"admin"
```

```json
"u0000"
```

```json
"worker$"
```

```json
"test-user"
```

```json
"test_user"
```

## set

Attributes of a User

`set`

* is required

* Type: `object` ([User attributes](cluster-definitions-user-attributes.md))

* cannot be null

* defined in: [alter-user input](cluster-definitions-user-attributes.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/set")

### set Type

`object` ([User attributes](cluster-definitions-user-attributes.md))

### set Examples

```json
{
  "display_name": "Mighty Admin",
  "2fa": false
}
```

## grant

A list of initial roles on the matching objects

`grant`

* is optional

* Type: `object[]` ([Grant object](cluster-definitions-grant-object.md))

* cannot be null

* defined in: [alter-user input](alter-user-input-properties-grant-assertions-list.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/grant")

### grant Type

`object[]` ([Grant object](cluster-definitions-grant-object.md))

### grant Constraints

**unique items**: all items in this array must be unique. Duplicates are not allowed.

## revoke

A list of revoke roles on the matching objects

`revoke`

* is optional

* Type: `object[]` ([Grant object](cluster-definitions-grant-object.md))

* cannot be null

* defined in: [alter-user input](alter-user-input-properties-revoke-assertions-list.md "http://schema.nethserver.org/cluster/alter-user-input.json#/properties/revoke")

### revoke Type

`object[]` ([Grant object](cluster-definitions-grant-object.md))

### revoke Constraints

**unique items**: all items in this array must be unique. Duplicates are not allowed.
