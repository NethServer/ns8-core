# Untitled object in list-backup-repositories output Schema

```txt
http://schema.nethserver.org/module/list-backup-repositories-output.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                   |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-backup-repositories-output.json\*](module/list-backup-repositories-output.json "open original schema") |

## items Type

`object` ([Details](list-backup-repositories-output-items.md))

# items Properties

| Property                                                     | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                |
| :----------------------------------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id)                                     | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-module_id.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/module_id")                                   |
| [module\_ui\_name](#module_ui_name)                          | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-module_ui_name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/module_ui_name")                         |
| [node\_fqdn](#node_fqdn)                                     | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-node_fqdn.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/node_fqdn")                                   |
| [path](#path)                                                | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-path.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/path")                                             |
| [name](#name)                                                | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/name")                                             |
| [uuid](#uuid)                                                | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-uuid.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/uuid")                                             |
| [timestamp](#timestamp)                                      | `integer` | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-timestamp.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/timestamp")                                   |
| [repository\_id](#repository_id)                             | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_id.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_id")                           |
| [repository\_name](#repository_name)                         | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_name")                       |
| [repository\_provider](#repository_provider)                 | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_provider.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_provider")               |
| [repository\_url](#repository_url)                           | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_url.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_url")                         |
| [installed\_instance](#installed_instance)                   | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-installed_instance.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/installed_instance")                 |
| [installed\_instance\_ui\_name](#installed_instance_ui_name) | `string`  | Required | cannot be null | [list-backup-repositories output](list-backup-repositories-output-items-properties-installed_instance_ui_name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/installed_instance_ui_name") |
| [is\_generated\_locally](#is_generated_locally)              | `boolean` | Required | can be null    | [list-backup-repositories output](list-backup-repositories-output-items-properties-is_generated_locally.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/is_generated_locally")             |

## module\_id

Original module ID value.

`module_id`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-module_id.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/module_id")

### module\_id Type

`string`

## module\_ui\_name

Original module label, assigned by the user.

`module_ui_name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-module_ui_name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/module_ui_name")

### module\_ui\_name Type

`string`

## node\_fqdn

The FQDN of the node where the module of the backup is hosted.

`node_fqdn`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-node_fqdn.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/node_fqdn")

### node\_fqdn Type

`string`

## path

Path of the repository, relative to the backup destination.

`path`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-path.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/path")

### path Type

`string`

## name

Name of the module. It is equal to the module image name.

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/name")

### name Type

`string`

## uuid

Universal, unique identifier of the module instance.

`uuid`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-uuid.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/uuid")

### uuid Type

`string`

## timestamp

Unix timestamp of the last backup run.

`timestamp`

* is required

* Type: `integer`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-timestamp.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/timestamp")

### timestamp Type

`integer`

## repository\_id

UUID of the backup destination.

`repository_id`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_id.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_id")

### repository\_id Type

`string`

## repository\_name

Human readable name of the backup destination.

`repository_name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_name")

### repository\_name Type

`string`

## repository\_provider

Type of backup destination provider, e.g. SMB, S3...

`repository_provider`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_provider.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_provider")

### repository\_provider Type

`string`

## repository\_url

Restic URL of the backup destination.

`repository_url`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-repository_url.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/repository_url")

### repository\_url Type

`string`

## installed\_instance

If the backup belongs to an installed module instance this is its module ID.

`installed_instance`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-installed_instance.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/installed_instance")

### installed\_instance Type

`string`

## installed\_instance\_ui\_name

If the backup belongs to an installed module instance this is its module friendly name.

`installed_instance_ui_name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-installed_instance_ui_name.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/installed_instance_ui_name")

### installed\_instance\_ui\_name Type

`string`

## is\_generated\_locally

Tells if the backup originates from the local cluster or from another cluster. The null value is returned if this information is missing completely, as it happens in old backups.

`is_generated_locally`

* is required

* Type: `boolean`

* can be null

* defined in: [list-backup-repositories output](list-backup-repositories-output-items-properties-is_generated_locally.md "http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/is_generated_locally")

### is\_generated\_locally Type

`boolean`
