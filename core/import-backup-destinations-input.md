# import-backup-destinations input Schema

```txt
http://schema.nethserver.org/cluster/import-backup-destinations-input.json
```

Import backup destination configurations from raw cluster backup data

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [import-backup-destinations-input.json](cluster/import-backup-destinations-input.json "open original schema") |

## import-backup-destinations input Type

`object` ([import-backup-destinations input](import-backup-destinations-input.md))

## import-backup-destinations input Examples

```json
{
  "backup_password": "mypassword",
  "backup_data": "jA0ECQMC2c...94FL/8y2KV"
}
```

# import-backup-destinations input Properties

| Property                            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                  |
| :---------------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [backup_password](#backup_password) | `string` | Required | cannot be null | [import-backup-destinations input](import-backup-destinations-input-properties-backup_password.md "http://schema.nethserver.org/cluster/import-backup-destinations-input.json#/properties/backup_password") |
| [backup_data](#backup_data)         | `string` | Required | cannot be null | [import-backup-destinations input](import-backup-destinations-input-properties-backup_data.md "http://schema.nethserver.org/cluster/import-backup-destinations-input.json#/properties/backup_data")         |

## backup_password

GPG passphrase to decrypt backup_data field

`backup_password`

* is required

* Type: `string`

* cannot be null

* defined in: [import-backup-destinations input](import-backup-destinations-input-properties-backup_password.md "http://schema.nethserver.org/cluster/import-backup-destinations-input.json#/properties/backup_password")

### backup_password Type

`string`

### backup_password Constraints

**maximum length**: the maximum number of characters for this string is: `4096`

## backup_data

base64-encoded string of GPG-encrypted, gzip-compressed backup data

`backup_data`

* is required

* Type: `string`

* cannot be null

* defined in: [import-backup-destinations input](import-backup-destinations-input-properties-backup_data.md "http://schema.nethserver.org/cluster/import-backup-destinations-input.json#/properties/backup_data")

### backup_data Type

`string`

### backup_data Constraints

**encoding**: the string content must be using the base64 content encoding.

**media type**: the media type of the contents of this string is: `application/octet-stream`
