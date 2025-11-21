# determine-restore-eligibility input Schema

```txt
http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json
```

Input schema of the determine-restore-eligibility action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [determine-restore-eligibility-input.json](cluster/determine-restore-eligibility-input.json "open original schema") |

## determine-restore-eligibility input Type

`object` ([determine-restore-eligibility input](determine-restore-eligibility-input.md))

## determine-restore-eligibility input Examples

```json
{
  "repository": "48ce000a-79b7-5fe6-8558-177fd70c27b4",
  "path": "dokuwiki/dokuwiki1@f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf",
  "snapshot": "a6b8317eef"
}
```

# determine-restore-eligibility input Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                       |
| :------------------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [repository](#repository) | `string` | Required | cannot be null | [determine-restore-eligibility input](determine-restore-eligibility-input-properties-destination-id.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json#/properties/repository")   |
| [path](#path)             | `string` | Required | cannot be null | [determine-restore-eligibility input](determine-restore-eligibility-input-properties-backup-repository-path.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json#/properties/path") |
| [snapshot](#snapshot)     | `string` | Required | cannot be null | [determine-restore-eligibility input](determine-restore-eligibility-input-properties-restic-snapshot-id.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json#/properties/snapshot") |

## repository

Backup destination identifier

`repository`

* is required

* Type: `string` ([Destination ID](determine-restore-eligibility-input-properties-destination-id.md))

* cannot be null

* defined in: [determine-restore-eligibility input](determine-restore-eligibility-input-properties-destination-id.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json#/properties/repository")

### repository Type

`string` ([Destination ID](determine-restore-eligibility-input-properties-destination-id.md))

### repository Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## path

Path of the Restic backup repository in the destination

`path`

* is required

* Type: `string` ([Backup repository path](determine-restore-eligibility-input-properties-backup-repository-path.md))

* cannot be null

* defined in: [determine-restore-eligibility input](determine-restore-eligibility-input-properties-backup-repository-path.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json#/properties/path")

### path Type

`string` ([Backup repository path](determine-restore-eligibility-input-properties-backup-repository-path.md))

### path Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## snapshot



`snapshot`

* is required

* Type: `string` ([Restic snapshot ID](determine-restore-eligibility-input-properties-restic-snapshot-id.md))

* cannot be null

* defined in: [determine-restore-eligibility input](determine-restore-eligibility-input-properties-restic-snapshot-id.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-input.json#/properties/snapshot")

### snapshot Type

`string` ([Restic snapshot ID](determine-restore-eligibility-input-properties-restic-snapshot-id.md))
