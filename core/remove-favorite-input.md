# remove-favorite input Schema

```txt
http://schema.nethserver.org/cluster/remove-favorite-input.json
```

Input schema of the remove-favorite action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-favorite-input.json](cluster/remove-favorite-input.json "open original schema") |

## remove-favorite input Type

`object` ([remove-favorite input](remove-favorite-input.md))

## remove-favorite input Examples

```json
{
  "instance": "dokuwiki1"
}
```

# remove-favorite input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                   |
| :-------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [instance](#instance) | `string` | Required | cannot be null | [remove-favorite input](remove-favorite-input-properties-instance.md "http://schema.nethserver.org/cluster/remove-favorite-input.json#/properties/instance") |

## instance

Instance name

`instance`

* is required

* Type: `string`

* cannot be null

* defined in: [remove-favorite input](remove-favorite-input-properties-instance.md "http://schema.nethserver.org/cluster/remove-favorite-input.json#/properties/instance")

### instance Type

`string`
