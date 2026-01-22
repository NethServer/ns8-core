# Azure blob storage protocol parameters Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## azure\_parameters Type

`object` ([Azure blob storage protocol parameters](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters.md))

# azure\_parameters Properties

| Property                                    | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                            |
| :------------------------------------------ | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [azure\_account\_name](#azure_account_name) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name") |
| [azure\_account\_key](#azure_account_key)   | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")   |

## azure\_account\_name



`azure_account_name`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name")

### azure\_account\_name Type

`string`

## azure\_account\_key



`azure_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")

### azure\_account\_key Type

`string`
