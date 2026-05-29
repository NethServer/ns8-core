# set-label input Schema

```txt
http://schema.nethserver.org/agent/set-label-input.json
```

Assign a user-defined name to the module instance. The name is shown in the UI.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-label-input.json](agent/set-label-input.json "open original schema") |

## set-label input Type

`object` ([set-label input](set-label-input.md))

## set-label input Examples

```json
{
  "name": "Nextcloud Milan"
}
```

# set-label input Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                       |
| :------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name) | `string` | Required | cannot be null | [set-label input](set-label-input-properties-name.md "http://schema.nethserver.org/agent/set-label-input.json#/properties/name") |

## name



`name`

* is required

* Type: `string`

* cannot be null

* defined in: [set-label input](set-label-input-properties-name.md "http://schema.nethserver.org/agent/set-label-input.json#/properties/name")

### name Type

`string`

### name Constraints

**maximum length**: the maximum number of characters for this string is: `24`
