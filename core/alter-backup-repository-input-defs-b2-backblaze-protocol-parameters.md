# B2 (Backblaze) protocol parameters Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## b2_parameters Type

`object` ([B2 (Backblaze) protocol parameters](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters.md))

# b2_parameters Properties

| Property                          | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                 |
| :-------------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [b2_account_id](#b2_account_id)   | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")   |
| [b2_account_key](#b2_account_key) | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key") |

## b2_account_id



`b2_account_id`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")

### b2_account_id Type

`string`

## b2_account_key



`b2_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key")

### b2_account_key Type

`string`
