# disable-node-terminal input Schema

```txt
http://schema.nethserver.org/cluster/disable-node-terminal-input.json
```

Forbid the cluster-admin terminal on a node. Closes live sessions and removes the sshd drop-in.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [disable-node-terminal-input.json](cluster/disable-node-terminal-input.json "open original schema") |

## disable-node-terminal input Type

`object` ([disable-node-terminal input](disable-node-terminal-input.md))

## disable-node-terminal input Examples

```json
{
  "node_id": 2
}
```

# disable-node-terminal input Properties

| Property             | Type      | Required | Nullable       | Defined by                                                                                                                                                                           |
| :------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id) | `integer` | Required | cannot be null | [disable-node-terminal input](disable-node-terminal-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/disable-node-terminal-input.json#/properties/node_id") |

## node\_id



`node_id`

* is required

* Type: `integer` ([Node identifier](disable-node-terminal-input-properties-node-identifier.md))

* cannot be null

* defined in: [disable-node-terminal input](disable-node-terminal-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/disable-node-terminal-input.json#/properties/node_id")

### node\_id Type

`integer` ([Node identifier](disable-node-terminal-input-properties-node-identifier.md))

### node\_id Constraints

**minimum**: the value of this number must greater than or equal to: `1`
