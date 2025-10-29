# S3-based provider schema Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/1/allOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## 0 Type

unknown ([S3-based provider schema](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema.md))

# 0 Properties

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                                |
| :------------------------ | :------------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [url](#url)               | Not specified | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-url.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/1/allOf/0/properties/url")               |
| [provider](#provider)     | Not specified | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/1/allOf/0/properties/provider") |
| [parameters](#parameters) | `object`      | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/1/allOf/0/properties/parameters")                       |

## url



`url`

* is optional

* Type: unknown

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-url.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/1/allOf/0/properties/url")

### url Type

unknown

### url Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
/
```

[try pattern](https://regexr.com/?expression=%2F "try regular expression with regexr.com")

## provider



`provider`

* is optional

* Type: unknown ([S3 providers](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/1/allOf/0/properties/provider")

### provider Type

unknown ([S3 providers](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md))

### provider Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value          | Explanation |
| :------------- | :---------- |
| `"aws"`        |             |
| `"generic-s3"` |             |

## parameters



`parameters`

* is optional

* Type: `object` ([S3 (Amazon AWS) protocol parameters](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/1/allOf/0/properties/parameters")

### parameters Type

`object` ([S3 (Amazon AWS) protocol parameters](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md))
