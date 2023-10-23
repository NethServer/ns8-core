# Untitled object in restore-module input Schema

```txt
http://schema.nethserver.org/cluster/restore-module-input.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [restore-module-input.json\*](cluster/restore-module-input.json "open original schema") |

## items Type

`object` ([Details](restore-module-input-1-items.md))

# items Properties

| Property                  | Type      | Required | Nullable       | Defined by                                                                                                                                                                     |
| :------------------------ | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [repository](#repository) | `string`  | Required | cannot be null | [restore-module input](restore-module-input-1-items-properties-repository-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/repository") |
| [path](#path)             | `string`  | Required | cannot be null | [restore-module input](restore-module-input-1-items-properties-backup-path.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/path")         |
| [snapshot](#snapshot)     | `string`  | Optional | cannot be null | [restore-module input](restore-module-input-1-items-properties-snapshot-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/snapshot")     |
| [node](#node)             | `integer` | Optional | cannot be null | [restore-module input](restore-module-input-1-items-properties-node-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/node")             |

## repository

Backup source repository

`repository`

*   is required

*   Type: `string` ([Repository ID](restore-module-input-1-items-properties-repository-id.md))

*   cannot be null

*   defined in: [restore-module input](restore-module-input-1-items-properties-repository-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/repository")

### repository Type

`string` ([Repository ID](restore-module-input-1-items-properties-repository-id.md))

### repository Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## path

Path of the backup in the source repository

`path`

*   is required

*   Type: `string` ([Backup path](restore-module-input-1-items-properties-backup-path.md))

*   cannot be null

*   defined in: [restore-module input](restore-module-input-1-items-properties-backup-path.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/path")

### path Type

`string` ([Backup path](restore-module-input-1-items-properties-backup-path.md))

### path Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## snapshot



`snapshot`

*   is optional

*   Type: `string` ([Snapshot ID](restore-module-input-1-items-properties-snapshot-id.md))

*   cannot be null

*   defined in: [restore-module input](restore-module-input-1-items-properties-snapshot-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/snapshot")

### snapshot Type

`string` ([Snapshot ID](restore-module-input-1-items-properties-snapshot-id.md))

## node

Node identifier where the module has to be restored

`node`

*   is optional

*   Type: `integer` ([Node ID](restore-module-input-1-items-properties-node-id.md))

*   cannot be null

*   defined in: [restore-module input](restore-module-input-1-items-properties-node-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/items/properties/node")

### node Type

`integer` ([Node ID](restore-module-input-1-items-properties-node-id.md))

### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`
