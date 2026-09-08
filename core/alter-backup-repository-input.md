# alter-backup-repository input Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json
```

Input schema of the alter-backup-repository action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json](cluster/alter-backup-repository-input.json "open original schema") |

## alter-backup-repository input Type

`object` ([alter-backup-repository input](alter-backup-repository-input.md))

any of

* all of

  * [Rclone schema](alter-backup-repository-input-anyof-0-allof-rclone-schema.md "check type definition")

* all of

  * [Safe values for parameters object](alter-backup-repository-input-anyof-1-allof-safe-values-for-parameters-object.md "check type definition")

  * [BackBlaze schema](alter-backup-repository-input-anyof-1-allof-backblaze-schema.md "check type definition")

* all of

  * [Safe values for parameters object](alter-backup-repository-input-anyof-2-allof-safe-values-for-parameters-object.md "check type definition")

  * [S3-based provider schema](alter-backup-repository-input-anyof-2-allof-s3-based-provider-schema.md "check type definition")

* all of

  * [Safe values for parameters object](alter-backup-repository-input-anyof-3-allof-safe-values-for-parameters-object.md "check type definition")

  * [Azure schema](alter-backup-repository-input-anyof-3-allof-azure-schema.md "check type definition")

* all of

  * [Safe values for parameters object](alter-backup-repository-input-anyof-4-allof-safe-values-for-parameters-object.md "check type definition")

  * [SMB schema](alter-backup-repository-input-anyof-4-allof-smb-schema.md "check type definition")

* all of

  * [Safe values for parameters object](alter-backup-repository-input-anyof-5-allof-safe-values-for-parameters-object.md "check type definition")

  * [Local storage](alter-backup-repository-input-anyof-5-allof-local-storage.md "check type definition")

## alter-backup-repository input Examples

```json
{
  "id:": "ce712dfe-b1d0-47f1-a149-a14d1203aab8",
  "provider": "backblaze",
  "name": "repository1",
  "parameters": {
    "b2_account_id": "YEFYEIOEHGO",
    "b2_account_key": "sdifjsiodv7sdv7suivyhsfv7fvyhdfvb7d"
  }
}
```

```json
{
  "name": "repository 6",
  "provider": "cluster",
  "url": "webdav:http://10.5.4.2:4694",
  "parameters": {}
}
```

```json
{
  "name": "destination 7",
  "provider": "rclone",
  "parameters": {
    "rclone_conf_secret": "type = s3\naccess_key_id = 45dc5e09346648fa9fed2a0267f58708\nsecret_access_key = db6b04c924f249fdbd86938b19959be0\nendpoint = s3.example.org\nprovider = Other\n",
    "basepath": "mybucket/subfolder"
  }
}
```

```json
{
  "name": "destination 7 - change path",
  "provider": "rclone",
  "parameters": {
    "rclone_conf_secret": "",
    "basepath": "mybucket/subfolder2"
  }
}
```

# alter-backup-repository input Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                          |
| :------------------------ | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)                 | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-properties-backup-repository-id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/id")          |
| [provider](#provider)     | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-properties-repository-provider.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/provider")     |
| [name](#name)             | `string` | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-properties-backup-repository-name.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/name")      |
| [parameters](#parameters) | `object` | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-properties-connection-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/parameters") |

## id



`id`

* is required

* Type: `string` ([Backup repository ID](alter-backup-repository-input-properties-backup-repository-id.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-properties-backup-repository-id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/id")

### id Type

`string` ([Backup repository ID](alter-backup-repository-input-properties-backup-repository-id.md))

### id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## provider



`provider`

* is required

* Type: `string` ([Repository provider](alter-backup-repository-input-properties-repository-provider.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-properties-repository-provider.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/provider")

### provider Type

`string` ([Repository provider](alter-backup-repository-input-properties-repository-provider.md))

### provider Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## name



`name`

* is optional

* Type: `string` ([Backup repository name](alter-backup-repository-input-properties-backup-repository-name.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-properties-backup-repository-name.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/name")

### name Type

`string` ([Backup repository name](alter-backup-repository-input-properties-backup-repository-name.md))

## parameters



`parameters`

* is optional

* Type: `object` ([Connection parameters](alter-backup-repository-input-properties-connection-parameters.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-properties-connection-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/properties/parameters")

### parameters Type

`object` ([Connection parameters](alter-backup-repository-input-properties-connection-parameters.md))

# alter-backup-repository input Definitions

## Definitions group safe\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/safe_parameters"}
```

| Property                    | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                                            |
| :-------------------------- | :------------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [parameters](#parameters-1) | Not specified | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-safe-values-for-parameters-object-properties-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/safe_parameters/properties/parameters") |

### parameters



`parameters`

* is optional

* Type: unknown

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-safe-values-for-parameters-object-properties-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/safe_parameters/properties/parameters")

#### parameters Type

unknown

## Definitions group rclone\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters"}
```

| Property                                             | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                             |
| :--------------------------------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [rclone\_conf\_secret](#rclone_conf_secret)          | `string` | Optional | can be null    | [alter-backup-repository input](alter-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_secret.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_secret")         |
| [rclone\_conf\_b64\_secret](#rclone_conf_b64_secret) | `string` | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_b64_secret.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_b64_secret") |
| [basepath](#basepath)                                | `string` | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-rclone-raw-configuration-properties-basepath.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters/properties/basepath")                             |

### rclone\_conf\_secret

Basic INI-like string validation

`rclone_conf_secret`

* is optional

* Type: `string`

* can be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_secret.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_secret")

#### rclone\_conf\_secret Type

`string`

#### rclone\_conf\_secret Constraints

**maximum length**: the maximum number of characters for this string is: `100000`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^(\[[^\]]+\]\n)?(\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*[^\n]*(\n|$))*$
```

[try pattern](https://regexr.com/?expression=%5E\(%5C%5B%5B%5E%5C%5D%5D%2B%5C%5D%5Cn\)%3F\(%5Cs*\(%5BA-Za-z_%5D%5BA-Za-z0-9_%5D*\)%5Cs*%3D%5Cs*%5B%5E%5Cn%5D*\(%5Cn%7C%24\)\)*%24 "try regular expression with regexr.com")

### rclone\_conf\_b64\_secret



`rclone_conf_b64_secret`

* is optional

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_b64_secret.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_b64_secret")

#### rclone\_conf\_b64\_secret Type

`string`

#### rclone\_conf\_b64\_secret Constraints

**maximum length**: the maximum number of characters for this string is: `100000`

### basepath



`basepath`

* is optional

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-rclone-raw-configuration-properties-basepath.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters/properties/basepath")

#### basepath Type

`string`

#### basepath Constraints

**maximum length**: the maximum number of characters for this string is: `512`

## Definitions group b2\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters"}
```

| Property                            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                 |
| :---------------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [b2\_account\_id](#b2_account_id)   | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")   |
| [b2\_account\_key](#b2_account_key) | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key") |

### b2\_account\_id



`b2_account_id`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")

#### b2\_account\_id Type

`string`

### b2\_account\_key



`b2_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key")

#### b2\_account\_key Type

`string`

## Definitions group s3\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters"}
```

| Property                                           | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                |
| :------------------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [aws\_default\_region](#aws_default_region)        | `string` | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")       |
| [aws\_access\_key\_id](#aws_access_key_id)         | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")         |
| [aws\_secret\_access\_key](#aws_secret_access_key) | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key") |

### aws\_default\_region



`aws_default_region`

* is optional

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")

#### aws\_default\_region Type

`string`

### aws\_access\_key\_id



`aws_access_key_id`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")

#### aws\_access\_key\_id Type

`string`

### aws\_secret\_access\_key



`aws_secret_access_key`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key")

#### aws\_secret\_access\_key Type

`string`

## Definitions group azure\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/azure_parameters"}
```

| Property                                    | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                  |
| :------------------------------------------ | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [azure\_account\_name](#azure_account_name) | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name") |
| [azure\_account\_key](#azure_account_key)   | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")   |

### azure\_account\_name



`azure_account_name`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name")

#### azure\_account\_name Type

`string`

### azure\_account\_key



`azure_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")

#### azure\_account\_key Type

`string`

## Definitions group smb\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters"}
```

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                               |
| :------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [smb\_host](#smb_host)     | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")     |
| [smb\_user](#smb_user)     | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")     |
| [smb\_pass](#smb_pass)     | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")     |
| [smb\_domain](#smb_domain) | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain") |

### smb\_host

Host name or IP address

`smb_host`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")

#### smb\_host Type

`string`

#### smb\_host Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")

### smb\_user

User name for share connection

`smb_user`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")

#### smb\_user Type

`string`

#### smb\_user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### smb\_pass

User password for share connection

`smb_pass`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")

#### smb\_pass Type

`string`

### smb\_domain

The short form (NT-style) domain name

`smb_domain`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain")

#### smb\_domain Type

`string`

## Definitions group cluster\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/cluster_parameters"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |
