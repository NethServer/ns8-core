# Azure schema Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/2/allOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## 0 Type

unknown ([Azure schema](add-backup-repository-input-anyof-2-allof-azure-schema.md))

# 0 Properties

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                      |
| :------------------------ | :------------ | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [provider](#provider)     | Not specified | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-anyof-2-allof-azure-schema-properties-azure-provider.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/2/allOf/0/properties/provider") |
| [parameters](#parameters) | `object`      | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/2/allOf/0/properties/parameters")        |

## provider



`provider`

* is optional

* Type: unknown ([Azure provider](add-backup-repository-input-anyof-2-allof-azure-schema-properties-azure-provider.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-anyof-2-allof-azure-schema-properties-azure-provider.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/2/allOf/0/properties/provider")

### provider Type

unknown ([Azure provider](add-backup-repository-input-anyof-2-allof-azure-schema-properties-azure-provider.md))

### provider Constraints

**constant**: the value of this property must be equal to:

```json
"azure"
```

## parameters



`parameters`

* is optional

* Type: `object` ([Azure blob storage protocol parameters](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/2/allOf/0/properties/parameters")

### parameters Type

`object` ([Azure blob storage protocol parameters](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters.md))
