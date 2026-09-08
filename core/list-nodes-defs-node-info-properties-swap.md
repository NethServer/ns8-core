# Untitled object in list-nodes output Schema

```txt
http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-nodes.json\*](cluster/list-nodes.json "open original schema") |

## swap Type

`object` ([Details](list-nodes-defs-node-info-properties-swap.md))

# swap Properties

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                                                                  |
| :-------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [total](#total) | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-swap-properties-total.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap/properties/total") |
| [used](#used)   | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-swap-properties-used.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap/properties/used")   |
| [free](#free)   | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-swap-properties-free.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap/properties/free")   |

## total



`total`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-swap-properties-total.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap/properties/total")

### total Type

`number`

## used



`used`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-swap-properties-used.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap/properties/used")

### used Type

`number`

## free



`free`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-swap-properties-free.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap/properties/free")

### free Type

`number`
