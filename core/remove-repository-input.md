# remove-repository input Schema

```txt
http://schema.nethserver.org/cluster/remove-repository-input.json
```

Input schema of the remove-repository action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-repository-input.json](cluster/remove-repository-input.json "open original schema") |

## remove-repository input Type

`object` ([remove-repository input](remove-repository-input.md))

## remove-repository input Examples

```json
{
  "name": "repository1"
}
```

# remove-repository input Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                 |
| :------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name) | `string` | Required | cannot be null | [remove-repository input](remove-repository-input-properties-name.md "http://schema.nethserver.org/cluster/remove-repository-input.json#/properties/name") |

## name

Repository name

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [remove-repository input](remove-repository-input-properties-name.md "http://schema.nethserver.org/cluster/remove-repository-input.json#/properties/name")

### name Type

`string`
