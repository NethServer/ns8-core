# Untitled object in Get module status Schema

```txt
http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-status-output.json\*](agent/get-status-output.json "open original schema") |

## items Type

`object` ([Details](get-status-output-properties-images-items.md))

# items Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                   |
| :------------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)       | `string` | Optional | cannot be null | [Get module status](get-status-output-properties-images-items-properties-name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items/properties/name")       |
| [size](#size)       | `string` | Optional | cannot be null | [Get module status](get-status-output-properties-images-items-properties-size.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items/properties/size")       |
| [created](#created) | `string` | Optional | cannot be null | [Get module status](get-status-output-properties-images-items-properties-created.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items/properties/created") |

## name

Name of image including repository and tag

`name`

* is optional

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-images-items-properties-name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items/properties/name")

### name Type

`string`

## size

The image size in human-readable format

`size`

* is optional

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-images-items-properties-size.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items/properties/size")

### size Type

`string`

## created

The image creation date

`created`

* is optional

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-images-items-properties-created.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images/items/properties/created")

### created Type

`string`
