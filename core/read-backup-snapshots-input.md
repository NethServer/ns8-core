# read-backup-snaphots input Schema

```txt
http://schema.nethserver.org/cluster/read-backup-snapshots-input.json
```

Input schema of the read-backup-snapshots action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [read-backup-snapshots-input.json](cluster/read-backup-snapshots-input.json "open original schema") |

## read-backup-snaphots input Type

`object` ([read-backup-snaphots input](read-backup-snapshots-input.md))

## read-backup-snaphots input Examples

```json
{
  "repository": "48ce000a-79b7-5fe6-8558-177fd70c27b4",
  "path": "dokuwiki/dokuwiki1@f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf"
}
```

# read-backup-snaphots input Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                           |
| :------------------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [repository](#repository) | `string` | Required | cannot be null | [read-backup-snaphots input](read-backup-snapshots-input-properties-repository-id.md "http://schema.nethserver.org/cluster/read-backup-snapshots-input.json#/properties/repository") |
| [path](#path)             | `string` | Required | cannot be null | [read-backup-snaphots input](read-backup-snapshots-input-properties-backup-path.md "http://schema.nethserver.org/cluster/read-backup-snapshots-input.json#/properties/path")         |

## repository

Backup source repository

`repository`

* is required

* Type: `string` ([Repository ID](read-backup-snapshots-input-properties-repository-id.md))

* cannot be null

* defined in: [read-backup-snaphots input](read-backup-snapshots-input-properties-repository-id.md "http://schema.nethserver.org/cluster/read-backup-snapshots-input.json#/properties/repository")

### repository Type

`string` ([Repository ID](read-backup-snapshots-input-properties-repository-id.md))

### repository Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## path

Path of the backup in the source repository

`path`

* is required

* Type: `string` ([Backup path](read-backup-snapshots-input-properties-backup-path.md))

* cannot be null

* defined in: [read-backup-snaphots input](read-backup-snapshots-input-properties-backup-path.md "http://schema.nethserver.org/cluster/read-backup-snapshots-input.json#/properties/path")

### path Type

`string` ([Backup path](read-backup-snapshots-input-properties-backup-path.md))

### path Constraints

**minimum length**: the minimum number of characters for this string is: `1`
