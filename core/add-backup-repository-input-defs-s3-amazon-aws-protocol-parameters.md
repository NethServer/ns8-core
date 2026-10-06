# S3 (Amazon AWS) protocol parameters Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## s3_parameters Type

`object` ([S3 (Amazon AWS) protocol parameters](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md))

# s3_parameters Properties

| Property                                        | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                          |
| :---------------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [aws_default_region](#aws_default_region)       | `string` | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")       |
| [aws_access_key_id](#aws_access_key_id)         | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")         |
| [aws_secret_access_key](#aws_secret_access_key) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key") |

## aws_default_region



`aws_default_region`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")

### aws_default_region Type

`string`

## aws_access_key_id



`aws_access_key_id`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")

### aws_access_key_id Type

`string`

## aws_secret_access_key



`aws_secret_access_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key")

### aws_secret_access_key Type

`string`
