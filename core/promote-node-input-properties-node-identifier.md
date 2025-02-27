# Node identifier Schema

```txt
http://schema.nethserver.org/cluster/promote-node-input.json#/properties/node_id
```

The node ID of the new leader node

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [promote-node-input.json\*](cluster/promote-node-input.json "open original schema") |

## node\_id Type

`integer` ([Node identifier](promote-node-input-properties-node-identifier.md))

## node\_id Constraints

**minimum**: the value of this number must greater than or equal to: `1`
