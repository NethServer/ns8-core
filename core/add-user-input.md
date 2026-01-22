# add-user input Schema

```txt
http://schema.nethserver.org/cluster/add-user-input.json
```

Create a user account in the Redis DB for the cluster administration web interface

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-user-input.json](cluster/add-user-input.json "open original schema") |

## add-user input Type

`object` ([add-user input](add-user-input.md))

## add-user input Examples

```json
{
  "user": "admin",
  "password_hash": "73cb3858a687a8494ca3323053016282f3dad39d42cf62ca4e79dda2aac7d9ac",
  "set": {
    "display_name": "The Administrator"
  },
  "grant": [
    {
      "role": "owner",
      "on": "*"
    }
  ]
}
```

# add-user input Properties

| Property                         | Type     | Required | Nullable       | Defined by                                                                                                                                        |
| :------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| [user](#user)                    | `string` | Required | cannot be null | [add-user input](cluster-definitions-strict-username-format.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/user")       |
| [password\_hash](#password_hash) | `string` | Required | cannot be null | [add-user input](cluster-definitions-redis-password.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/password_hash")      |
| [set](#set)                      | `object` | Required | cannot be null | [add-user input](cluster-definitions-user-attributes.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/set")               |
| [grant](#grant)                  | `array`  | Required | cannot be null | [add-user input](add-user-input-properties-grant-assertions-list.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/grant") |

## user

A username string can be written in many ways. This is a quite strict form. See <https://www.unix.com/man-page/linux/8/useradd/>

`user`

* is required

* Type: `string` ([Strict username format](cluster-definitions-strict-username-format.md))

* cannot be null

* defined in: [add-user input](cluster-definitions-strict-username-format.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/user")

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

## password\_hash

Redis SHA256 hashed password for Users, Cluster, Nodes and Modules

`password_hash`

* is required

* Type: `string` ([Redis password](cluster-definitions-redis-password.md))

* cannot be null

* defined in: [add-user input](cluster-definitions-redis-password.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/password_hash")

### password\_hash Type

`string` ([Redis password](cluster-definitions-redis-password.md))

### password\_hash Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `64`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[a-f0-9]{64}$
```

[try pattern](https://regexr.com/?expression=%5E%5Ba-f0-9%5D%7B64%7D%24 "try regular expression with regexr.com")

### password\_hash Examples

```json
"73cb3858a687a8494ca3323053016282f3dad39d42cf62ca4e79dda2aac7d9ac"
```

## set

Attributes of a User

`set`

* is required

* Type: `object` ([User attributes](cluster-definitions-user-attributes.md))

* cannot be null

* defined in: [add-user input](cluster-definitions-user-attributes.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/set")

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

* is required

* Type: `object[]` ([Grant object](cluster-definitions-grant-object.md))

* cannot be null

* defined in: [add-user input](add-user-input-properties-grant-assertions-list.md "http://schema.nethserver.org/cluster/add-user-input.json#/properties/grant")

### grant Type

`object[]` ([Grant object](cluster-definitions-grant-object.md))

### grant Constraints

**unique items**: all items in this array must be unique. Duplicates are not allowed.
