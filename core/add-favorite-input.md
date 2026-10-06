# add-favorite input Schema

```txt
http://schema.nethserver.org/cluster/add-favorite-input.json
```

Input schema of the add-favorite action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-favorite-input.json](cluster/add-favorite-input.json "open original schema") |

## add-favorite input Type

`object` ([add-favorite input](add-favorite-input.md))

## add-favorite input Examples

```json
{
  "instance": "dokuwiki1"
}
```

# add-favorite input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                          |
| :-------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| [instance](#instance) | `string` | Required | cannot be null | [add-favorite input](add-favorite-input-properties-instance.md "http://schema.nethserver.org/cluster/add-favorite-input.json#/properties/instance") |

## instance

Instance name

`instance`

* is required

* Type: `string`

* cannot be null

* defined in: [add-favorite input](add-favorite-input-properties-instance.md "http://schema.nethserver.org/cluster/add-favorite-input.json#/properties/instance")

### instance Type

`string`
