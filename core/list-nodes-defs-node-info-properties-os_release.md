# Untitled object in list-nodes output Schema

```txt
http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/os_release
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-nodes.json\*](cluster/list-nodes.json "open original schema") |

## os\_release Type

`object` ([Details](list-nodes-defs-node-info-properties-os_release.md))

# os\_release Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                  |
| :------------------ | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)       | `string` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-os_release-properties-name.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/os_release/properties/name")       |
| [version](#version) | `string` | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-os_release-properties-version.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/os_release/properties/version") |

## name



`name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-os_release-properties-name.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/os_release/properties/name")

### name Type

`string`

### name Examples

```json
"Rocky Linux"
```

## version



`version`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-os_release-properties-version.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/os_release/properties/version")

### version Type

`string`

### version Examples

```json
"9.6"
```
