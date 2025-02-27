# User attributes Schema

```txt
http://schema.nethserver.org/cluster/alter-user-input.json#/properties/set
```

Attributes of a User

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-user-input.json\*](cluster/alter-user-input.json "open original schema") |

## set Type

`object` ([User attributes](cluster-definitions-user-attributes.md))

## set Examples

```json
{
  "display_name": "Mighty Admin",
  "2fa": false
}
```

# set Properties

| Property                       | Type      | Required | Nullable       | Defined by                                                                                                                                                                         |
| :----------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [display\_name](#display_name) | `string`  | Optional | cannot be null | [Cluster library](cluster-definitions-user-attributes-properties-display_name.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/display_name") |
| [2fa](#2fa)                    | `boolean` | Optional | cannot be null | [Cluster library](cluster-definitions-user-attributes-properties-2fa.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/2fa")                   |

## display\_name



`display_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [Cluster library](cluster-definitions-user-attributes-properties-display_name.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/display_name")

### display\_name Type

`string`

### display\_name Constraints

**minimum length**: the minimum number of characters for this string is: `0`

## 2fa



`2fa`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [Cluster library](cluster-definitions-user-attributes-properties-2fa.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/2fa")

### 2fa Type

`boolean`
