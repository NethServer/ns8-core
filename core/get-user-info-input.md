# get-user-info input Schema

```txt
http://schema.nethserver.org/cluster/get-user-info-input.json
```

Input schema of the get-user-info action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-user-info-input.json](cluster/get-user-info-input.json "open original schema") |

## get-user-info input Type

`object` ([get-user-info input](get-user-info-input.md))

## get-user-info input Examples

```json
{
  "user": "admin"
}
```

# get-user-info input Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                     |
| :------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------- |
| [user](#user) | `string` | Required | cannot be null | [get-user-info input](get-user-info-input-properties-user.md "http://schema.nethserver.org/cluster/get-user-info-input.json#/properties/user") |

## user

Target user.

`user`

* is required

* Type: `string`

* cannot be null

* defined in: [get-user-info input](get-user-info-input-properties-user.md "http://schema.nethserver.org/cluster/get-user-info-input.json#/properties/user")

### user Type

`string`
