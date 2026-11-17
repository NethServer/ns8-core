# Backup file Schema

```txt
http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/file
```

Backup file encoded with base64

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [retrieve-cluster-backup-input.json\*](cluster/retrieve-cluster-backup-input.json "open original schema") |

## file Type

`string` ([Backup file](retrieve-cluster-backup-input-properties-backup-file.md))

## file Constraints

**encoding**: the string content must be using the base64 content encoding.

**media type**: the media type of the contents of this string is: `application/octet-stream`
