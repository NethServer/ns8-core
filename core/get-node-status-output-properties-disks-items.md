# Untitled object in get-node-status output Schema

```txt
http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-node-status-output.json\*](node/get-node-status-output.json "open original schema") |

## items Type

`object` ([Details](get-node-status-output-properties-disks-items.md))

# items Properties

| Property                  | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                     |
| :------------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [total](#total)           | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-disks-items-properties-total.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/total")           |
| [used](#used)             | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-disks-items-properties-used.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/used")             |
| [free](#free)             | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-disks-items-properties-free.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/free")             |
| [device](#device)         | `string`  | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-disks-items-properties-device.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/device")         |
| [mountpoint](#mountpoint) | `string`  | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-disks-items-properties-mountpoint.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/mountpoint") |
| [fstype](#fstype)         | `string`  | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-disks-items-properties-fstype.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/fstype")         |

## total

Total space

`total`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-disks-items-properties-total.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/total")

### total Type

`integer`

## used

Used space

`used`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-disks-items-properties-used.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/used")

### used Type

`integer`

## free

Free space

`free`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-disks-items-properties-free.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/free")

### free Type

`integer`

## device

Partition device path

`device`

* is optional

* Type: `string`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-disks-items-properties-device.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/device")

### device Type

`string`

## mountpoint

Mount point

`mountpoint`

* is optional

* Type: `string`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-disks-items-properties-mountpoint.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/mountpoint")

### mountpoint Type

`string`

## fstype

File system type

`fstype`

* is optional

* Type: `string`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-disks-items-properties-fstype.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks/items/properties/fstype")

### fstype Type

`string`
