# Untitled object in get-node-status output Schema

```txt
http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-node-status-output.json\*](node/get-node-status-output.json "open original schema") |

## items Type

`object` ([Details](get-node-status-output-properties-cpu-properties-info-items.md))

# items Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                         |
| :---------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [vendor](#vendor) | `string` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-cpu-properties-info-items-properties-vendor.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items/properties/vendor") |
| [model](#model)   | `string` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-cpu-properties-info-items-properties-model.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items/properties/model")   |
| [clock](#clock)   | `string` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-cpu-properties-info-items-properties-clock.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items/properties/clock")   |

## vendor

CPU vendor

`vendor`

* is optional

* Type: `string`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-cpu-properties-info-items-properties-vendor.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items/properties/vendor")

### vendor Type

`string`

## model

CPU model

`model`

* is optional

* Type: `string`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-cpu-properties-info-items-properties-model.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items/properties/model")

### model Type

`string`

## clock

CPU frequency in MHz

`clock`

* is optional

* Type: `string`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-cpu-properties-info-items-properties-clock.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info/items/properties/clock")

### clock Type

`string`
