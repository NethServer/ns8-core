# BackBlaze schema Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/0/allOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## 0 Type

unknown ([BackBlaze schema](alter-backup-repository-input-anyof-0-allof-backblaze-schema.md))

# 0 Properties

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                                    |
| :------------------------ | :------------ | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [provider](#provider)     | Not specified | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-anyof-0-allof-backblaze-schema-properties-backblaze-provider.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/0/allOf/0/properties/provider") |
| [parameters](#parameters) | `object`      | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/0/allOf/0/properties/parameters")                      |

## provider



`provider`

* is optional

* Type: unknown ([BackBlaze provider](alter-backup-repository-input-anyof-0-allof-backblaze-schema-properties-backblaze-provider.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-anyof-0-allof-backblaze-schema-properties-backblaze-provider.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/0/allOf/0/properties/provider")

### provider Type

unknown ([BackBlaze provider](alter-backup-repository-input-anyof-0-allof-backblaze-schema-properties-backblaze-provider.md))

### provider Constraints

**constant**: the value of this property must be equal to:

```json
"backblaze"
```

## parameters



`parameters`

* is optional

* Type: `object` ([B2 (Backblaze) protocol parameters](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/0/allOf/0/properties/parameters")

### parameters Type

`object` ([B2 (Backblaze) protocol parameters](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters.md))
