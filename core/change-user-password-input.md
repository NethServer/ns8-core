# change-user-password-input Schema

```txt
http://schema.nethserver.org/cluster/change-user-password-input.json
```

Input schema of the change-user-password action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [change-user-password-input.json](cluster/change-user-password-input.json "open original schema") |

## change-user-password-input Type

`object` ([change-user-password-input](change-user-password-input.md))

## change-user-password-input Examples

```json
{
  "user": "admin",
  "current_password": "Nethesis,1234",
  "new_password": "Nethesis,12345"
}
```

# change-user-password-input Properties

| Property                               | Type     | Required | Nullable       | Defined by                                                                                                                                                                                  |
| :------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [user](#user)                          | `string` | Required | cannot be null | [change-user-password-input](change-user-password-input-properties-user.md "http://schema.nethserver.org/cluster/change-user-password-input.json#/properties/user")                         |
| [current\_password](#current_password) | `string` | Required | cannot be null | [change-user-password-input](change-user-password-input-properties-current_password.md "http://schema.nethserver.org/cluster/change-user-password-input.json#/properties/current_password") |
| [new\_password](#new_password)         | `string` | Required | cannot be null | [change-user-password-input](change-user-password-input-properties-new_password.md "http://schema.nethserver.org/cluster/change-user-password-input.json#/properties/new_password")         |

## user

Target user.

`user`

* is required

* Type: `string`

* cannot be null

* defined in: [change-user-password-input](change-user-password-input-properties-user.md "http://schema.nethserver.org/cluster/change-user-password-input.json#/properties/user")

### user Type

`string`

## current\_password

Current password in plain text

`current_password`

* is required

* Type: `string`

* cannot be null

* defined in: [change-user-password-input](change-user-password-input-properties-current_password.md "http://schema.nethserver.org/cluster/change-user-password-input.json#/properties/current_password")

### current\_password Type

`string`

## new\_password

New password in plain text. The password will be updated only if the authentication with 'user' and 'current\_password' is valid.

`new_password`

* is required

* Type: `string`

* cannot be null

* defined in: [change-user-password-input](change-user-password-input-properties-new_password.md "http://schema.nethserver.org/cluster/change-user-password-input.json#/properties/new_password")

### new\_password Type

`string`
