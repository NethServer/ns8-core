# Untitled object in list-core-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-core-modules-output.json\*](cluster/list-core-modules-output.json "open original schema") |

## items Type

`object` ([Details](list-core-modules-output-items-properties-instances-items.md))

# items Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                            |
| :------------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [id](#id)           | `string` | Required | cannot be null | [list-core-modules output](list-core-modules-output-items-properties-instances-items-properties-id.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items/properties/id")           |
| [version](#version) | `string` | Required | cannot be null | [list-core-modules output](list-core-modules-output-items-properties-instances-items-properties-version.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items/properties/version") |
| [update](#update)   | `string` | Required | cannot be null | [list-core-modules output](list-core-modules-output-items-properties-instances-items-properties-update.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items/properties/update")   |

## id

Instance id

`id`

*   is required

*   Type: `string`

*   cannot be null

*   defined in: [list-core-modules output](list-core-modules-output-items-properties-instances-items-properties-id.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items/properties/id")

### id Type

`string`

## version

Installed version

`version`

*   is required

*   Type: `string`

*   cannot be null

*   defined in: [list-core-modules output](list-core-modules-output-items-properties-instances-items-properties-version.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items/properties/version")

### version Type

`string`

## update

Available version update, can be empty

`update`

*   is required

*   Type: `string`

*   cannot be null

*   defined in: [list-core-modules output](list-core-modules-output-items-properties-instances-items-properties-update.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances/items/properties/update")

### update Type

`string`
