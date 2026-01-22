# get-user-info output Schema

```txt
http://schema.nethserver.org/cluster/get-user-info-output.json
```

Output schema of the get-user-info action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-user-info-output.json](cluster/get-user-info-output.json "open original schema") |

## get-user-info output Type

`object` ([get-user-info output](get-user-info-output.md))

## get-user-info output Examples

```json
{
  "display_name": "Admin",
  "email": ""
}
```

# get-user-info output Properties

| Property                       | Type     | Required | Nullable       | Defined by                                                                                                                                                        |
| :----------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [display\_name](#display_name) | `string` | Required | cannot be null | [get-user-info output](get-user-info-output-properties-display_name.md "http://schema.nethserver.org/cluster/get-user-info-output.json#/properties/display_name") |
| [email](#email)                | `string` | Required | cannot be null | [get-user-info output](get-user-info-output-properties-email.md "http://schema.nethserver.org/cluster/get-user-info-output.json#/properties/email")               |

## display\_name

Display name of the user

`display_name`

* is required

* Type: `string`

* cannot be null

* defined in: [get-user-info output](get-user-info-output-properties-display_name.md "http://schema.nethserver.org/cluster/get-user-info-output.json#/properties/display_name")

### display\_name Type

`string`

## email

Email of the user

`email`

* is required

* Type: `string`

* cannot be null

* defined in: [get-user-info output](get-user-info-output-properties-email.md "http://schema.nethserver.org/cluster/get-user-info-output.json#/properties/email")

### email Type

`string`
