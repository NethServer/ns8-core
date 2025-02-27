# join-cluster output Schema

```txt
http://schema.nethserver.org/cluster/join-cluster-output.json
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [join-cluster-output.json](cluster/join-cluster-output.json "open original schema") |

## join-cluster output Type

`object` ([join-cluster output](join-cluster-output.md))

## join-cluster output Examples

```json
{
  "node_id": 3
}
```

# join-cluster output Properties

| Property             | Type      | Required | Nullable       | Defined by                                                                                                                                                   |
| :------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id) | `integer` | Required | cannot be null | [join-cluster output](join-cluster-output-properties-node-identifier.md "http://schema.nethserver.org/cluster/join-cluster-output.json#/properties/node_id") |

## node\_id

Allocated identifier of the new node

`node_id`

* is required

* Type: `integer` ([Node identifier](join-cluster-output-properties-node-identifier.md))

* cannot be null

* defined in: [join-cluster output](join-cluster-output-properties-node-identifier.md "http://schema.nethserver.org/cluster/join-cluster-output.json#/properties/node_id")

### node\_id Type

`integer` ([Node identifier](join-cluster-output-properties-node-identifier.md))
