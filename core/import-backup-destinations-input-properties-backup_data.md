# Untitled string in import-backup-destinations input Schema

```txt
http://schema.nethserver.org/cluster/import-backup-destinations-input.json#/properties/backup_data
```

base64-encoded string of GPG-encrypted, gzip-compressed backup data

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [import-backup-destinations-input.json\*](cluster/import-backup-destinations-input.json "open original schema") |

## backup\_data Type

`string`

## backup\_data Constraints

**encoding**: the string content must be using the base64 content encoding.

**media type**: the media type of the contents of this string is: `application/octet-stream`
