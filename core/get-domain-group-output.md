# get-domain-group output Schema

```txt
http://schema.nethserver.org/cluster/get-domain-group-output.json
```

Return the details of a group

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-domain-group-output.json](cluster/get-domain-group-output.json "open original schema") |

## get-domain-group output Type

`object` ([get-domain-group output](get-domain-group-output.md))

## get-domain-group output Examples

```json
{
  "group": {
    "group": "Domain Admins",
    "description": "Designated administrators of the domain",
    "users": [
      {
        "user": "Administrator",
        "display_name": "Administrator"
      }
    ]
  }
}
```

# get-domain-group output Properties

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                                     |
| :-------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [group](#group) | `object` | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-group-details.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/properties/group") |

## group



`group`

* is required

* Type: `object` ([Group details](get-domain-group-output-defs-group-details.md))

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-group-details.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/properties/group")

### group Type

`object` ([Group details](get-domain-group-output-defs-group-details.md))

# get-domain-group output Definitions

## Definitions group user

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user"}
```

| Property                       | Type     | Required | Nullable       | Defined by                                                                                                                                                                                      |
| :----------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [user](#user)                  | `string` | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-user-properties-user.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/user")                 |
| [display\_name](#display_name) | `string` | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-user-properties-display-name.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/display_name") |

### user



`user`

* is required

* Type: `string`

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-user-properties-user.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/user")

#### user Type

`string`

#### user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### display\_name

Preferred name of a person to be used when displaying entries (RFC 2798/2.3)

`display_name`

* is required

* Type: `string` ([Display name](get-domain-group-output-defs-user-properties-display-name.md))

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-user-properties-display-name.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/user/properties/display_name")

#### display\_name Type

`string` ([Display name](get-domain-group-output-defs-user-properties-display-name.md))

## Definitions group group

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group"}
```

| Property                    | Type     | Required | Nullable       | Defined by                                                                                                                                                                                              |
| :-------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [group](#group-1)           | `string` | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-group-details-properties-group.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group/properties/group")             |
| [description](#description) | `string` | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-group-details-properties-description.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group/properties/description") |
| [users](#users)             | `array`  | Required | cannot be null | [get-domain-group output](get-domain-group-output-defs-group-details-properties-group-members.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group/properties/users")     |

### group



`group`

* is required

* Type: `string`

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-group-details-properties-group.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group/properties/group")

#### group Type

`string`

#### group Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### description

A brief description of the group purpose

`description`

* is required

* Type: `string` ([Description](get-domain-group-output-defs-group-details-properties-description.md))

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-group-details-properties-description.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group/properties/description")

#### description Type

`string` ([Description](get-domain-group-output-defs-group-details-properties-description.md))

### users



`users`

* is required

* Type: `object[]` ([Details](get-domain-group-output-defs-user.md))

* cannot be null

* defined in: [get-domain-group output](get-domain-group-output-defs-group-details-properties-group-members.md "http://schema.nethserver.org/cluster/get-domain-group-output.json#/$defs/group/properties/users")

#### users Type

`object[]` ([Details](get-domain-group-output-defs-user.md))
