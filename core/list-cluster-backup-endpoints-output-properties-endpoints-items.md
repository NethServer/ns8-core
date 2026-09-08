# Untitled object in list-cluster-backup-endpoints output Schema

```txt
http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json#/properties/endpoints/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-cluster-backup-endpoints-output.json\*](cluster/list-cluster-backup-endpoints-output.json "open original schema") |

## items Type

`object` ([Details](list-cluster-backup-endpoints-output-properties-endpoints-items.md))

# items Properties

| Property               | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                      |
| :--------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ui\_label](#ui_label) | `string` | Required | cannot be null | [list-cluster-backup-endpoints output](list-cluster-backup-endpoints-output-properties-endpoints-items-properties-ui_label.md "http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json#/properties/endpoints/items/properties/ui_label") |
| [url](#url)            | `string` | Required | cannot be null | [list-cluster-backup-endpoints output](list-cluster-backup-endpoints-output-properties-endpoints-items-properties-url.md "http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json#/properties/endpoints/items/properties/url")           |

## ui\_label

Display label

`ui_label`

* is required

* Type: `string`

* cannot be null

* defined in: [list-cluster-backup-endpoints output](list-cluster-backup-endpoints-output-properties-endpoints-items-properties-ui_label.md "http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json#/properties/endpoints/items/properties/ui_label")

### ui\_label Type

`string`

## url

Base URL of rclone-gateway WebDAV service

`url`

* is required

* Type: `string`

* cannot be null

* defined in: [list-cluster-backup-endpoints output](list-cluster-backup-endpoints-output-properties-endpoints-items-properties-url.md "http://schema.nethserver.org/cluster/list-cluster-backup-endpoints-output.json#/properties/endpoints/items/properties/url")

### url Type

`string`
