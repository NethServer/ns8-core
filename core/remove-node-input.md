# remove-node input Schema

```txt
http://schema.nethserver.org/cluster/remove-node-input.json
```

Remove a node from the cluster

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-node-input.json](cluster/remove-node-input.json "open original schema") |

## remove-node input Type

`object` ([remove-node input](remove-node-input.md))

## remove-node input Examples

```json
{
  "node_id": 3
}
```

# remove-node input Properties

| Property             | Type      | Required | Nullable       | Defined by                                                                                                                                             |
| :------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id) | `integer` | Required | cannot be null | [remove-node input](remove-node-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/remove-node-input.json#/properties/node_id") |

## node\_id



`node_id`

* is required

* Type: `integer` ([Node identifier](remove-node-input-properties-node-identifier.md))

* cannot be null

* defined in: [remove-node input](remove-node-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/remove-node-input.json#/properties/node_id")

### node\_id Type

`integer` ([Node identifier](remove-node-input-properties-node-identifier.md))

### node\_id Constraints

**minimum**: the value of this number must greater than or equal to: `1`
