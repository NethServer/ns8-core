# Untitled object in list-updates output Schema

```txt
http://schema.nethserver.org/cluster/list-updates-output.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-updates-output.json\*](cluster/list-updates-output.json "open original schema") |

## items Type

`object` ([Details](list-updates-output-items.md))

# items Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                       |
| :------------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)           | `string` | Required | cannot be null | [list-updates output](list-updates-output-items-properties-id.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/id")           |
| [node](#node)       | `string` | Required | cannot be null | [list-updates output](list-updates-output-items-properties-node.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/node")       |
| [digest](#digest)   | `string` | Required | cannot be null | [list-updates output](list-updates-output-items-properties-digest.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/digest")   |
| [source](#source)   | `string` | Required | cannot be null | [list-updates output](list-updates-output-items-properties-source.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/source")   |
| [version](#version) | `string` | Required | cannot be null | [list-updates output](list-updates-output-items-properties-version.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/version") |
| [update](#update)   | `string` | Required | cannot be null | [list-updates output](list-updates-output-items-properties-update.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/update")   |

## id

Unique name of a module instance

`id`

* is required

* Type: `string`

* cannot be null

* defined in: [list-updates output](list-updates-output-items-properties-id.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/id")

### id Type

`string`

## node

Id of the node where the instance is running

`node`

* is required

* Type: `string`

* cannot be null

* defined in: [list-updates output](list-updates-output-items-properties-node.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/node")

### node Type

`string`

## digest

Image digest

`digest`

* is required

* Type: `string`

* cannot be null

* defined in: [list-updates output](list-updates-output-items-properties-digest.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/digest")

### digest Type

`string`

## source

The URL of the container image registry

`source`

* is required

* Type: `string`

* cannot be null

* defined in: [list-updates output](list-updates-output-items-properties-source.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/source")

### source Type

`string`

## version

A valid semantic version extracted from image tag

`version`

* is required

* Type: `string`

* cannot be null

* defined in: [list-updates output](list-updates-output-items-properties-version.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/version")

### version Type

`string`

## update

A valid semantic version extracted from image tag wich should be greater than 'version' field

`update`

* is required

* Type: `string`

* cannot be null

* defined in: [list-updates output](list-updates-output-items-properties-update.md "http://schema.nethserver.org/cluster/list-updates-output.json#/items/properties/update")

### update Type

`string`
