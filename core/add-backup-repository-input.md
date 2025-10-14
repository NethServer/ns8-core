# add-backup-repository input Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json
```

Input schema of the add-backup-repository action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-backup-repository-input.json](cluster/add-backup-repository-input.json "open original schema") |

## add-backup-repository input Type

`object` ([add-backup-repository input](add-backup-repository-input.md))

any of

* all of

  * [BackBlaze schema](add-backup-repository-input-anyof-0-allof-backblaze-schema.md "check type definition")

* all of

  * [S3-based provider schema](add-backup-repository-input-anyof-1-allof-s3-based-provider-schema.md "check type definition")

* all of

  * [Azure schema](add-backup-repository-input-anyof-2-allof-azure-schema.md "check type definition")

* all of

  * [SMB schema](add-backup-repository-input-anyof-3-allof-smb-schema.md "check type definition")

* all of

  * [Local storage](add-backup-repository-input-anyof-4-allof-local-storage.md "check type definition")

## add-backup-repository input Examples

```json
{
  "name": "repository1",
  "provider": "backblaze",
  "url": "b2:mybucket/mybackup",
  "password": "",
  "parameters": {
    "b2_account_id": "YEFYEIOEHGO",
    "b2_account_key": "sdifjsiodv7sdv7suivyhsfv7fvyhdfvb7d"
  }
}
```

```json
{
  "name": "repository2",
  "provider": "aws",
  "url": "s3:s3.amazonaws.com/mybucket/mybackup",
  "password": "45a1905ef8de3c03b05d47071754bd5ddbfec0edaa56d4c44f981386f3f24888",
  "parameters": {
    "aws_default_region": "us-east-1",
    "aws_access_key_id": "IEIEEHEHEW",
    "aws_secret_access_key": "edfjksof798r7fsdfiougvf7df"
  }
}
```

```json
{
  "name": "repository 4",
  "provider": "azure",
  "password": "",
  "url": "azure:mystorage",
  "parameters": {
    "azure_account_name": "ns8backup",
    "azure_account_key": "7lBkSMdy5LeYFyjcc0Szlxlc9GDaUFxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxOrplWQJuEBRD+AStiQnppg=="
  }
}
```

```json
{
  "name": "repository 5",
  "provider": "smb",
  "password": "",
  "url": "rclone::smb:myshare/subpath",
  "parameters": {
    "rclone_smb_host": "10.114.0.2",
    "rclone_smb_user": "backup-agent",
    "rclone_smb_pass": "Nethesis,1234",
    "rclone_smb_domain": "NTDOM"
  }
}
```

```json
{
  "name": "repository 6",
  "provider": "cluster",
  "password": "",
  "url": "webdav:http://10.5.4.2:4694",
  "parameters": {}
}
```

# add-backup-repository input Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                    |
| :------------------------ | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [provider](#provider)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-properties-repository-provider.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/provider")     |
| [name](#name)             | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-properties-backup-repository-name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/name")      |
| [url](#url)               | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-properties-restic-url.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/url")                   |
| [password](#password)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-properties-encryption-token.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/password")        |
| [parameters](#parameters) | `object` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-properties-connection-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/parameters") |

## provider



`provider`

* is required

* Type: `string` ([Repository provider](add-backup-repository-input-properties-repository-provider.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-properties-repository-provider.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/provider")

### provider Type

`string` ([Repository provider](add-backup-repository-input-properties-repository-provider.md))

### provider Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## name



`name`

* is required

* Type: `string` ([Backup repository name](add-backup-repository-input-properties-backup-repository-name.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-properties-backup-repository-name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/name")

### name Type

`string` ([Backup repository name](add-backup-repository-input-properties-backup-repository-name.md))

## url

URL of the repository in restic format. Must be unique.

`url`

* is required

* Type: `string` ([Restic URL](add-backup-repository-input-properties-restic-url.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-properties-restic-url.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/url")

### url Type

`string` ([Restic URL](add-backup-repository-input-properties-restic-url.md))

## password

Select the password for restic encryption. If this is empty the API will generate a random password

`password`

* is required

* Type: `string` ([Encryption token](add-backup-repository-input-properties-encryption-token.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-properties-encryption-token.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/password")

### password Type

`string` ([Encryption token](add-backup-repository-input-properties-encryption-token.md))

## parameters



`parameters`

* is required

* Type: `object` ([Connection parameters](add-backup-repository-input-properties-connection-parameters.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-properties-connection-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/parameters")

### parameters Type

`object` ([Connection parameters](add-backup-repository-input-properties-connection-parameters.md))

# add-backup-repository input Definitions

## Definitions group b2\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters"}
```

| Property                            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                           |
| :---------------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [b2\_account\_id](#b2_account_id)   | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")   |
| [b2\_account\_key](#b2_account_key) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key") |

### b2\_account\_id



`b2_account_id`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")

#### b2\_account\_id Type

`string`

### b2\_account\_key



`b2_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key")

#### b2\_account\_key Type

`string`

## Definitions group s3\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters"}
```

| Property                                           | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                          |
| :------------------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [aws\_default\_region](#aws_default_region)        | `string` | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")       |
| [aws\_access\_key\_id](#aws_access_key_id)         | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")         |
| [aws\_secret\_access\_key](#aws_secret_access_key) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key") |

### aws\_default\_region



`aws_default_region`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")

#### aws\_default\_region Type

`string`

### aws\_access\_key\_id



`aws_access_key_id`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")

#### aws\_access\_key\_id Type

`string`

### aws\_secret\_access\_key



`aws_secret_access_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key")

#### aws\_secret\_access\_key Type

`string`

## Definitions group azure\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters"}
```

| Property                                    | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                            |
| :------------------------------------------ | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [azure\_account\_name](#azure_account_name) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name") |
| [azure\_account\_key](#azure_account_key)   | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")   |

### azure\_account\_name



`azure_account_name`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name")

#### azure\_account\_name Type

`string`

### azure\_account\_key



`azure_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")

#### azure\_account\_key Type

`string`

## Definitions group smb\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters"}
```

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                         |
| :------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [smb\_host](#smb_host)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")     |
| [smb\_user](#smb_user)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")     |
| [smb\_pass](#smb_pass)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")     |
| [smb\_domain](#smb_domain) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain") |

### smb\_host

Host name or IP address

`smb_host`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")

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

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")

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

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")

#### smb\_pass Type

`string`

### smb\_domain

The short form (NT-style) domain name

`smb_domain`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain")

#### smb\_domain Type

`string`

## Definitions group cluster\_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/cluster_parameters"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |
