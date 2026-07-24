# Untitled object in list-users output Schema

```txt
http://schema.nethserver.org/cluster/list-users-output.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-users-output.json\*](cluster/list-users-output.json "open original schema") |

## items Type

`object` ([Details](list-users-output-items.md))

# items Properties

| Property                       | Type     | Required | Nullable       | Defined by                                                                                                                                                           |
| :----------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [display\_name](#display_name) | `string` | Required | cannot be null | [list-users output](list-users-output-items-properties-display_name.md "http://schema.nethserver.org/cluster/list-users-output.json#/items/properties/display_name") |
| [user](#user)                  | `string` | Required | cannot be null | [list-users output](list-users-output-items-properties-user.md "http://schema.nethserver.org/cluster/list-users-output.json#/items/properties/user")                 |

## display\_name

Display name of the user

`display_name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-users output](list-users-output-items-properties-display_name.md "http://schema.nethserver.org/cluster/list-users-output.json#/items/properties/display_name")

### display\_name Type

`string`

## user

Unique username for the user

`user`

* is required

* Type: `string`

* cannot be null

* defined in: [list-users output](list-users-output-items-properties-user.md "http://schema.nethserver.org/cluster/list-users-output.json#/items/properties/user")

### user Type

`string`
