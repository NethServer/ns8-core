# Backup of a module instance Schema

```txt
http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-backups-output.json\*](cluster/list-backups-output.json "open original schema") |

## instance-item Type

`object` ([Backup of a module instance](list-backups-output-defs-backup-of-a-module-instance.md))

# instance-item Properties

| Property                             | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                |
| :----------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [module\_id](#module_id)             | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-module_id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/module_id")             |
| [ui\_name](#ui_name)                 | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-ui_name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/ui_name")                 |
| [repository\_path](#repository_path) | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-repository_path.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/repository_path") |
| [status](#status)                    | `object` | Optional | can be null    | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/status")   |

## module\_id



`module_id`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-module_id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/module_id")

### module\_id Type

`string`

### module\_id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## ui\_name



`ui_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-ui_name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/ui_name")

### ui\_name Type

`string`

## repository\_path



`repository_path`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-repository_path.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/repository_path")

### repository\_path Type

`string`

### repository\_path Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## status



`status`

* is optional

* Type: `object` ([Last backup run status](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md))

* can be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/status")

### status Type

`object` ([Last backup run status](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md))
