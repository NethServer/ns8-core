# Untitled object in list-nodes output Schema

```txt
http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-nodes.json\*](cluster/list-nodes.json "open original schema") |

## memory Type

`object` ([Details](list-nodes-defs-node-info-properties-memory.md))

# memory Properties

| Property        | Type     | Required | Nullable       | Defined by                                                                                                                                                                                      |
| :-------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [total](#total) | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-memory-properties-total.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory/properties/total") |
| [used](#used)   | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-memory-properties-used.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory/properties/used")   |
| [free](#free)   | `number` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-memory-properties-free.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory/properties/free")   |

## total



`total`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-memory-properties-total.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory/properties/total")

### total Type

`number`

## used



`used`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-memory-properties-used.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory/properties/used")

### used Type

`number`

## free



`free`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-memory-properties-free.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory/properties/free")

### free Type

`number`
