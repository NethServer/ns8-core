# Untitled object in list-nodes output Schema

```txt
http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu
```

Average load

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-nodes.json\*](cluster/list-nodes.json "open original schema") |

## cpu Type

`object` ([Details](list-nodes-defs-node-info-properties-cpu.md))

# cpu Properties

| Property                   | Type      | Required | Nullable       | Defined by                                                                                                                                                                                          |
| :------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [count](#count)            | `integer` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/count")           |
| [usage](#usage)            | `number`  | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-usage.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/usage")           |
| [vendor](#vendor)          | `string`  | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-vendor.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/vendor")         |
| [model](#model)            | `string`  | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-model.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/model")           |
| [model\_name](#model_name) | `string`  | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-model_name.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/model_name") |
| [family](#family)          | `string`  | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-family.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/family")         |
| [stepping](#stepping)      | `string`  | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-stepping.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/stepping")     |
| [package](#package)        | `string`  | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-package.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/package")       |

## count

Number of CPUs

`count`

* is optional

* Type: `integer`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/count")

### count Type

`integer`

## usage

Normalized value of CPU time in busy (not in idle) state

`usage`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-usage.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/usage")

### usage Type

`number`

## vendor



`vendor`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-vendor.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/vendor")

### vendor Type

`string`

## model



`model`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-model.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/model")

### model Type

`string`

## model\_name



`model_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-model_name.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/model_name")

### model\_name Type

`string`

## family



`family`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-family.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/family")

### family Type

`string`

## stepping



`stepping`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-stepping.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/stepping")

### stepping Type

`string`

## package



`package`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu-properties-package.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu/properties/package")

### package Type

`string`
