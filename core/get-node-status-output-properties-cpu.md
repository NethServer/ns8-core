# Untitled object in get-node-status output Schema

```txt
http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu
```

Average load

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-node-status-output.json\*](node/get-node-status-output.json "open original schema") |

## cpu Type

`object` ([Details](get-node-status-output-properties-cpu.md))

# cpu Properties

| Property        | Type      | Required | Nullable       | Defined by                                                                                                                                                                           |
| :-------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [count](#count) | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-cpu-properties-count.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/count") |
| [usage](#usage) | `number`  | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-cpu-properties-usage.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/usage") |
| [info](#info)   | `array`   | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-cpu-properties-info.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info")   |

## count

Number of physical and virtual CPUs

`count`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-cpu-properties-count.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/count")

### count Type

`integer`

## usage

Global CPU percentage usage

`usage`

* is optional

* Type: `number`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-cpu-properties-usage.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/usage")

### usage Type

`number`

## info

CPU details

`info`

* is optional

* Type: `object[]` ([Details](get-node-status-output-properties-cpu-properties-info-items.md))

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-cpu-properties-info.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu/properties/info")

### info Type

`object[]` ([Details](get-node-status-output-properties-cpu-properties-info-items.md))
