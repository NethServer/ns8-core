# list-cluster-backup-endpoints output Schema

```txt
http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json
```

Get the list of WebDAV backup endpoints provided by cluster nodes

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-cluster-backup-endpoints-output.json](cluster/list-cluster-backup-endpoints-output.json "open original schema") |

## list-cluster-backup-endpoints output Type

`object` ([list-cluster-backup-endpoints output](list-cluster-backup-endpoints-output.md))

## list-cluster-backup-endpoints output Examples

```json
{
  "endpoints": [
    {
      "ui_label": "Node 1",
      "url": "webdav:http://10.5.4.1:4694"
    }
  ]
}
```

# list-cluster-backup-endpoints output Properties

| Property                | Type    | Required | Nullable       | Defined by                                                                                                                                                                                                  |
| :---------------------- | :------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [endpoints](#endpoints) | `array` | Required | cannot be null | [list-cluster-backup-endpoints output](list-cluster-backup-endpoints-output-properties-endpoints.md "http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json#/properties/endpoints") |

## endpoints



`endpoints`

* is required

* Type: `object[]` ([Details](list-cluster-backup-endpoints-output-properties-endpoints-items.md))

* cannot be null

* defined in: [list-cluster-backup-endpoints output](list-cluster-backup-endpoints-output-properties-endpoints.md "http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json#/properties/endpoints")

### endpoints Type

`object[]` ([Details](list-cluster-backup-endpoints-output-properties-endpoints-items.md))
