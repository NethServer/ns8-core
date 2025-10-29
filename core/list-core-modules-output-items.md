# Untitled object in list-core-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-core-modules-output.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-core-modules-output.json\*](cluster/list-core-modules-output.json "open original schema") |

## items Type

`object` ([Details](list-core-modules-output-items.md))

# items Properties

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                                          |
| :---------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)           | `string` | Required | cannot be null | [list-core-modules output](list-core-modules-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/name")           |
| [instances](#instances) | `array`  | Required | cannot be null | [list-core-modules output](list-core-modules-output-items-properties-instances.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances") |

## name

Unique name of a package

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-core-modules output](list-core-modules-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/name")

### name Type

`string`

## instances



`instances`

* is required

* Type: `object[]` ([Details](list-core-modules-output-items-properties-instances-items.md))

* cannot be null

* defined in: [list-core-modules output](list-core-modules-output-items-properties-instances.md "http://schema.nethserver.org/cluster/list-core-modules-output.json#/items/properties/instances")

### instances Type

`object[]` ([Details](list-core-modules-output-items-properties-instances-items.md))
