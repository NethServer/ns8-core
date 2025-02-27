# Untitled object in get-node-status output Schema

```txt
http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory
```

Memory statistics

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-node-status-output.json\*](node/get-node-status-output.json "open original schema") |

## memory Type

`object` ([Details](get-node-status-output-properties-memory.md))

# memory Properties

| Property        | Type      | Required | Nullable       | Defined by                                                                                                                                                                                 |
| :-------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [total](#total) | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-memory-properties-total.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory/properties/total") |
| [used](#used)   | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-memory-properties-used.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory/properties/used")   |
| [free](#free)   | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-memory-properties-free.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory/properties/free")   |

## total

Total memory

`total`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-memory-properties-total.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory/properties/total")

### total Type

`integer`

## used

Used memory

`used`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-memory-properties-used.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory/properties/used")

### used Type

`integer`

## free

Free memory

`free`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-memory-properties-free.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory/properties/free")

### free Type

`integer`
