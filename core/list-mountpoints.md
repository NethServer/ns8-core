# list-mountpoints output Schema

```txt
http://schema.nethserver.org/node/list-mountpoints.json
```

Output schema of the list-mountpoints action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-mountpoints.json](node/list-mountpoints.json "open original schema") |

## list-mountpoints output Type

`object` ([list-mountpoints output](list-mountpoints.md))

## list-mountpoints output Examples

```json
{
  "default_disk": {
    "path": "/",
    "size": 63254278144,
    "used": 11392909312,
    "label": "rl1_root",
    "ui_name": "rl1_root",
    "available": 51861368832
  },
  "mountpoints": [
    {
      "ui_name": "disk00",
      "path": "/mnt/disk00",
      "label": "",
      "size": 2136997888,
      "used": 49049600
    },
    {
      "ui_name": "VOL01",
      "path": "/mnt/xvolume_davidep_cluster0_rl1",
      "label": "VOL01",
      "size": 2136997888,
      "used": 49049600
    }
  ]
}
```

# list-mountpoints output Properties

| Property                       | Type     | Required | Nullable       | Defined by                                                                                                                                              |
| :----------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [default\_disk](#default_disk) | `object` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo.md "http://schema.nethserver.org/node/list-mountpoints.json#/properties/default_disk")           |
| [mountpoints](#mountpoints)    | `array`  | Optional | cannot be null | [list-mountpoints output](list-mountpoints-properties-mountpoints.md "http://schema.nethserver.org/node/list-mountpoints.json#/properties/mountpoints") |

## default\_disk



`default_disk`

* is optional

* Type: `object` ([Details](list-mountpoints-defs-fsinfo.md))

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo.md "http://schema.nethserver.org/node/list-mountpoints.json#/properties/default_disk")

### default\_disk Type

`object` ([Details](list-mountpoints-defs-fsinfo.md))

## mountpoints



`mountpoints`

* is optional

* Type: `object[]` ([Details](list-mountpoints-defs-fsinfo.md))

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-properties-mountpoints.md "http://schema.nethserver.org/node/list-mountpoints.json#/properties/mountpoints")

### mountpoints Type

`object[]` ([Details](list-mountpoints-defs-fsinfo.md))

# list-mountpoints output Definitions

## Definitions group fsinfo

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo"}
```

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                                   |
| :---------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ui\_name](#ui_name)    | `string` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-ui_name.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/ui_name")     |
| [path](#path)           | `string` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-path.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/path")           |
| [label](#label)         | `string` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-label.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/label")         |
| [size](#size)           | `number` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-size.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/size")           |
| [available](#available) | `number` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-available.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/available") |
| [used](#used)           | `number` | Optional | cannot be null | [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-used.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/used")           |

### ui\_name



`ui_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-ui_name.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/ui_name")

#### ui\_name Type

`string`

### path



`path`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-path.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/path")

#### path Type

`string`

### label



`label`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-label.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/label")

#### label Type

`string`

### size

bytes

`size`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-size.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/size")

#### size Type

`number`

### available

bytes

`available`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-available.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/available")

#### available Type

`number`

### used

bytes

`used`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-mountpoints output](list-mountpoints-defs-fsinfo-properties-used.md "http://schema.nethserver.org/node/list-mountpoints.json#/$defs/fsinfo/properties/used")

#### used Type

`number`
