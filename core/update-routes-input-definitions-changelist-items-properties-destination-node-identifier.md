# Destination node identifier Schema

```txt
http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/node_id
```

Node ID used as route next-hop

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [update-routes-input.json\*](cluster/update-routes-input.json "open original schema") |

## node\_id Type

`integer` ([Destination node identifier](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md))

## node\_id Constraints

**minimum (exclusive)**: the value of this number must be greater than: `0`
