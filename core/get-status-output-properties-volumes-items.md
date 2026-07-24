# Untitled object in Get module status Schema

```txt
http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-status-output.json\*](agent/get-status-output.json "open original schema") |

## items Type

`object` ([Details](get-status-output-properties-volumes-items.md))

# items Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                     |
| :------------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)       | `string` | Optional | cannot be null | [Get module status](get-status-output-properties-volumes-items-properties-name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items/properties/name")       |
| [mount](#mount)     | `string` | Optional | cannot be null | [Get module status](get-status-output-properties-volumes-items-properties-mount.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items/properties/mount")     |
| [created](#created) | `string` | Optional | cannot be null | [Get module status](get-status-output-properties-volumes-items-properties-created.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items/properties/created") |

## name

Name of volume

`name`

* is optional

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-volumes-items-properties-name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items/properties/name")

### name Type

`string`

## mount

Mount point of the volume

`mount`

* is optional

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-volumes-items-properties-mount.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items/properties/mount")

### mount Type

`string`

## created

The volume creation date

`created`

* is optional

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-volumes-items-properties-created.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes/items/properties/created")

### created Type

`string`
