# S3-based provider schema Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/1/allOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## 0 Type

unknown ([S3-based provider schema](alter-backup-repository-input-anyof-1-allof-s3-based-provider-schema.md))

# 0 Properties

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                                      |
| :------------------------ | :------------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [provider](#provider)     | Not specified | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/1/allOf/0/properties/provider") |
| [parameters](#parameters) | `object`      | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/1/allOf/0/properties/parameters")                       |

## provider



`provider`

* is optional

* Type: unknown ([S3 providers](alter-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/1/allOf/0/properties/provider")

### provider Type

unknown ([S3 providers](alter-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md))

### provider Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value          | Explanation |
| :------------- | :---------- |
| `"aws"`        |             |
| `"generic-s3"` |             |

## parameters



`parameters`

* is optional

* Type: `object` ([S3 (Amazon AWS) protocol parameters](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/1/allOf/0/properties/parameters")

### parameters Type

`object` ([S3 (Amazon AWS) protocol parameters](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md))
