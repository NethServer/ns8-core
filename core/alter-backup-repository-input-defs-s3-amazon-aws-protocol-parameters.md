# S3 (Amazon AWS) protocol parameters Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## s3\_parameters Type

`object` ([S3 (Amazon AWS) protocol parameters](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters.md))

# s3\_parameters Properties

| Property                                           | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                |
| :------------------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [aws\_default\_region](#aws_default_region)        | `string` | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")       |
| [aws\_access\_key\_id](#aws_access_key_id)         | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")         |
| [aws\_secret\_access\_key](#aws_secret_access_key) | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key") |

## aws\_default\_region



`aws_default_region`

* is optional

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")

### aws\_default\_region Type

`string`

## aws\_access\_key\_id



`aws_access_key_id`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")

### aws\_access\_key\_id Type

`string`

## aws\_secret\_access\_key



`aws_secret_access_key`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key")

### aws\_secret\_access\_key Type

`string`
