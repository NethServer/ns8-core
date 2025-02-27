# Untitled object in get-node-status output Schema

```txt
http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap
```

SWAP statistics

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-node-status-output.json\*](node/get-node-status-output.json "open original schema") |

## swap Type

`object` ([Details](get-node-status-output-properties-swap.md))

# swap Properties

| Property        | Type      | Required | Nullable       | Defined by                                                                                                                                                                             |
| :-------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [total](#total) | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-swap-properties-total.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap/properties/total") |
| [used](#used)   | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-swap-properties-used.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap/properties/used")   |
| [free](#free)   | `integer` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-swap-properties-free.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap/properties/free")   |

## total

Total SWAP

`total`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-swap-properties-total.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap/properties/total")

### total Type

`integer`

## used

Used SWAP

`used`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-swap-properties-used.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap/properties/used")

### used Type

`integer`

## free

Free SWAP

`free`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-swap-properties-free.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap/properties/free")

### free Type

`integer`
