# restore-module input Schema

```txt
http://schema.nethserver.org/cluster/restore-module-input.json
```

Input schema of the restore-module action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [restore-module-input.json](cluster/restore-module-input.json "open original schema") |

## restore-module input Type

`object` ([restore-module input](restore-module-input-1.md))

## restore-module input Examples

```json
{
  "repository": "48ce000a-79b7-5fe6-8558-177fd70c27b4",
  "path": "dokuwiki/dokuwiki1@f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf",
  "snapshot": "",
  "node": 1,
  "replace": false
}
```

# restore-module input Properties

| Property                  | Type      | Required | Nullable       | Defined by                                                                                                                                                         |
| :------------------------ | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [repository](#repository) | `string`  | Required | cannot be null | [restore-module input](restore-module-input-1-properties-repository-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/repository") |
| [path](#path)             | `string`  | Required | cannot be null | [restore-module input](restore-module-input-1-properties-backup-path.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/path")         |
| [snapshot](#snapshot)     | `string`  | Required | cannot be null | [restore-module input](restore-module-input-1-properties-snapshot-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/snapshot")     |
| [node](#node)             | `integer` | Required | cannot be null | [restore-module input](restore-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/node")             |
| [replace](#replace)       | `boolean` | Optional | cannot be null | [restore-module input](restore-module-input-1-properties-replace-flag.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/replace")     |

## repository

Backup source repository

`repository`

* is required

* Type: `string` ([Repository ID](restore-module-input-1-properties-repository-id.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-1-properties-repository-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/repository")

### repository Type

`string` ([Repository ID](restore-module-input-1-properties-repository-id.md))

### repository Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## path

Path of the backup in the source repository

`path`

* is required

* Type: `string` ([Backup path](restore-module-input-1-properties-backup-path.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-1-properties-backup-path.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/path")

### path Type

`string` ([Backup path](restore-module-input-1-properties-backup-path.md))

### path Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## snapshot



`snapshot`

* is required

* Type: `string` ([Snapshot ID](restore-module-input-1-properties-snapshot-id.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-1-properties-snapshot-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/snapshot")

### snapshot Type

`string` ([Snapshot ID](restore-module-input-1-properties-snapshot-id.md))

## node

Node identifier where the module has to be restored

`node`

* is required

* Type: `integer` ([Node ID](restore-module-input-1-properties-node-id.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/node")

### node Type

`integer` ([Node ID](restore-module-input-1-properties-node-id.md))

### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## replace

If set to true, remove any existing instance with the same UUID, then restore the UUID to the new module

`replace`

* is optional

* Type: `boolean` ([Replace flag](restore-module-input-1-properties-replace-flag.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-1-properties-replace-flag.md "http://schema.nethserver.org/cluster/restore-module-input.json#/properties/replace")

### replace Type

`boolean` ([Replace flag](restore-module-input-1-properties-replace-flag.md))
