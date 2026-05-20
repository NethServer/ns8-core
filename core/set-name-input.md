# set-name input Schema

```txt
http://schema.nethserver.org/node/set-name-input.json
```

Assign a user-defined name to the node instance. The name is shown in the UI.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                             |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-name-input.json](node/set-name-input.json "open original schema") |

## set-name input Type

`object` ([set-name input](set-name-input.md))

## set-name input Examples

```json
{
  "name": "Milan0"
}
```

# set-name input Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                   |
| :------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------- |
| [name](#name) | `string` | Required | cannot be null | [set-name input](set-name-input-properties-name.md "http://schema.nethserver.org/node/set-name-input.json#/properties/name") |

## name



`name`

* is required

* Type: `string`

* cannot be null

* defined in: [set-name input](set-name-input-properties-name.md "http://schema.nethserver.org/node/set-name-input.json#/properties/name")

### name Type

`string`

### name Constraints

**maximum length**: the maximum number of characters for this string is: `24`
