# Untitled object in list-mountpoints output Schema

```txt
http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                   |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-mountpoints.json\*](node/list-mountpoints.json "open original schema") |

## fsinfo Type

`object` ([Details](list-mountpoints-defs-fsinfo.md))

# fsinfo Properties

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                                   |
| :---------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ui\_name](#ui_name)    | `string` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-ui_name.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/ui_name")     |
| [path](#path)           | `string` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-path.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/path")           |
| [label](#label)         | `string` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-label.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/label")         |
| [size](#size)           | `number` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-size.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/size")           |
| [available](#available) | `number` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-available.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/available") |
| [used](#used)           | `number` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-used.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/used")           |

## ui\_name



`ui_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-ui_name.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/ui_name")

### ui\_name Type

`string`

## path



`path`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-path.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/path")

### path Type

`string`

## label



`label`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-label.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/label")

### label Type

`string`

## size

bytes

`size`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-size.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/size")

### size Type

`number`

## available

bytes

`available`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-available.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/available")

### available Type

`number`

## used

bytes

`used`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-used.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/used")

### used Type

`number`
