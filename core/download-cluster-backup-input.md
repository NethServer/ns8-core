# download-cluster-backup-input Schema

```txt
http://schema.nethserver.org/cluster/download-cluster-backup-input.json
```

Input schema of the download-cluster-backup-input action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [download-cluster-backup-input.json](cluster/download-cluster-backup-input.json "open original schema") |

## download-cluster-backup-input Type

`object` ([download-cluster-backup-input](download-cluster-backup-input.md))

## download-cluster-backup-input Examples

```json
{
  "password": "Nethesis,12345"
}
```

# download-cluster-backup-input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                                           |
| :-------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [password](#password) | `string` | Optional | cannot be null | [download-cluster-backup-input](download-cluster-backup-input-properties-password.md "http://schema.nethserver.org/cluster/download-cluster-backup-input.json#/properties/password") |

## password

Encryption password in plain text

`password`

* is optional

* Type: `string`

* cannot be null

* defined in: [download-cluster-backup-input](download-cluster-backup-input-properties-password.md "http://schema.nethserver.org/cluster/download-cluster-backup-input.json#/properties/password")

### password Type

`string`

### password Constraints

**minimum length**: the minimum number of characters for this string is: `1`
