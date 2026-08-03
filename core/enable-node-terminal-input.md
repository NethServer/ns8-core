# enable-node-terminal input Schema

```txt
http://schema.nethserver.org/cluster/enable-node-terminal-input.json
```

Allow the cluster-admin terminal on a node. Installs an sshd drop-in that accepts password authentication from the cluster VPN only.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [enable-node-terminal-input.json](cluster/enable-node-terminal-input.json "open original schema") |

## enable-node-terminal input Type

`object` ([enable-node-terminal input](enable-node-terminal-input.md))

## enable-node-terminal input Examples

```json
{
  "node_id": 2
}
```

# enable-node-terminal input Properties

| Property             | Type      | Required | Nullable       | Defined by                                                                                                                                                                        |
| :------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id) | `integer` | Required | cannot be null | [enable-node-terminal input](enable-node-terminal-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/enable-node-terminal-input.json#/properties/node_id") |

## node\_id



`node_id`

* is required

* Type: `integer` ([Node identifier](enable-node-terminal-input-properties-node-identifier.md))

* cannot be null

* defined in: [enable-node-terminal input](enable-node-terminal-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/enable-node-terminal-input.json#/properties/node_id")

### node\_id Type

`integer` ([Node identifier](enable-node-terminal-input-properties-node-identifier.md))

### node\_id Constraints

**minimum**: the value of this number must greater than or equal to: `1`
