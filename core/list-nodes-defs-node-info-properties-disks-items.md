# Untitled object in list-nodes output Schema

```txt
http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-nodes.json\*](cluster/list-nodes.json "open original schema") |

## items Type

`object` ([Details](list-nodes-defs-node-info-properties-disks-items.md))

# items Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                          |
| :------------------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [total](#total)           | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-total.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/total")           |
| [used](#used)             | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-used.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/used")             |
| [free](#free)             | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-free.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/free")             |
| [device](#device)         | `string` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-device.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/device")         |
| [mountpoint](#mountpoint) | `string` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-mountpoint.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/mountpoint") |
| [fstype](#fstype)         | `string` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-fstype.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/fstype")         |

## total



`total`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-total.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/total")

### total Type

`number`

## used



`used`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-used.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/used")

### used Type

`number`

## free



`free`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-free.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/free")

### free Type

`number`

## device



`device`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-device.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/device")

### device Type

`string`

## mountpoint



`mountpoint`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-mountpoint.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/mountpoint")

### mountpoint Type

`string`

## fstype



`fstype`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-disks-items-properties-fstype.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks/items/properties/fstype")

### fstype Type

`string`

### fstype Examples

```json
"ext4"
```

```json
"xfs"
```
