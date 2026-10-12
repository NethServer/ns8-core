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

  * [Rclone schema](add-backup-repository-input-anyof-0-allof-rclone-schema.md "check type definition")

* all of

  * [Safe values for parameters object](add-backup-repository-input-anyof-1-allof-safe-values-for-parameters-object.md "check type definition")

  * [BackBlaze schema](add-backup-repository-input-anyof-1-allof-backblaze-schema.md "check type definition")

* all of

  * [Safe values for parameters object](add-backup-repository-input-anyof-2-allof-safe-values-for-parameters-object.md "check type definition")

  * [S3-based provider schema](add-backup-repository-input-anyof-2-allof-s3-based-provider-schema.md "check type definition")

* all of

  * [Safe values for parameters object](add-backup-repository-input-anyof-3-allof-safe-values-for-parameters-object.md "check type definition")

  * [Azure schema](add-backup-repository-input-anyof-3-allof-azure-schema.md "check type definition")

* all of

  * [Safe values for parameters object](add-backup-repository-input-anyof-4-allof-safe-values-for-parameters-object.md "check type definition")

  * [SMB schema](add-backup-repository-input-anyof-4-allof-smb-schema.md "check type definition")

* all of

  * [Safe values for parameters object](add-backup-repository-input-anyof-5-allof-safe-values-for-parameters-object.md "check type definition")

  * [Local storage](add-backup-repository-input-anyof-5-allof-local-storage.md "check type definition")

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

```json
{
  "name": "destination 7",
  "provider": "rclone",
  "password": "",
  "url": "",
  "parameters": {
    "rclone_conf_secret": "type = s3\naccess_key_id = 45dc5e09346648fa9fed2a0267f58708\nsecret_access_key = db6b04c924f249fdbd86938b19959be0\nendpoint = s3.example.org\nprovider = Other\n",
    "basepath": "mybucket/subfolder"
  }
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

## Definitions group safe_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/safe_parameters"}
```

| Property                    | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                                      |
| :-------------------------- | :------------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [parameters](#parameters-1) | Not specified | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-safe-values-for-parameters-object-properties-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/safe_parameters/properties/parameters") |

### parameters



`parameters`

* is optional

* Type: unknown

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-safe-values-for-parameters-object-properties-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/safe_parameters/properties/parameters")

#### parameters Type

unknown

## Definitions group rclone_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters"}
```

| Property                                          | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                       |
| :------------------------------------------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [rclone_conf_secret](#rclone_conf_secret)         | `string` | Required | can be null    | [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_secret")         |
| [rclone_conf_b64_secret](#rclone_conf_b64_secret) | `string` | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_b64_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_b64_secret") |
| [basepath](#basepath)                             | `string` | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-basepath.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/basepath")                             |

### rclone_conf_secret

Basic INI-like string validation

`rclone_conf_secret`

* is required

* Type: `string`

* can be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_secret")

#### rclone_conf_secret Type

`string`

#### rclone_conf_secret Constraints

**maximum length**: the maximum number of characters for this string is: `100000`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^(\[[^\]]+\]\n)?(\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*[^\n]*(\n|$))*$
```

[try pattern](https://regexr.com/?expression=%5E\(%5C%5B%5B%5E%5C%5D%5D%2B%5C%5D%5Cn\)%3F\(%5Cs*\(%5BA-Za-z_%5D%5BA-Za-z0-9_%5D*\)%5Cs*%3D%5Cs*%5B%5E%5Cn%5D*\(%5Cn%7C%24\)\)*%24 "try regular expression with regexr.com")

### rclone_conf_b64_secret



`rclone_conf_b64_secret`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_b64_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_b64_secret")

#### rclone_conf_b64_secret Type

`string`

#### rclone_conf_b64_secret Constraints

**maximum length**: the maximum number of characters for this string is: `100000`

### basepath



`basepath`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-basepath.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/basepath")

#### basepath Type

`string`

#### basepath Constraints

**maximum length**: the maximum number of characters for this string is: `512`

## Definitions group b2_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters"}
```

| Property                          | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                           |
| :-------------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [b2_account_id](#b2_account_id)   | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")   |
| [b2_account_key](#b2_account_key) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key") |

### b2_account_id



`b2_account_id`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_id")

#### b2_account_id Type

`string`

### b2_account_key



`b2_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-b2-backblaze-protocol-parameters-properties-b2_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/b2_parameters/properties/b2_account_key")

#### b2_account_key Type

`string`

## Definitions group s3_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters"}
```

| Property                                        | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                          |
| :---------------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [aws_default_region](#aws_default_region)       | `string` | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")       |
| [aws_access_key_id](#aws_access_key_id)         | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")         |
| [aws_secret_access_key](#aws_secret_access_key) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key") |

### aws_default_region



`aws_default_region`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_default_region.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_default_region")

#### aws_default_region Type

`string`

### aws_access_key_id



`aws_access_key_id`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_access_key_id.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_access_key_id")

#### aws_access_key_id Type

`string`

### aws_secret_access_key



`aws_secret_access_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-s3-amazon-aws-protocol-parameters-properties-aws_secret_access_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/s3_parameters/properties/aws_secret_access_key")

#### aws_secret_access_key Type

`string`

## Definitions group azure_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters"}
```

| Property                                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                            |
| :---------------------------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [azure_account_name](#azure_account_name) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name") |
| [azure_account_key](#azure_account_key)   | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")   |

### azure_account_name



`azure_account_name`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_name.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_name")

#### azure_account_name Type

`string`

### azure_account_key



`azure_account_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-azure-blob-storage-protocol-parameters-properties-azure_account_key.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/azure_parameters/properties/azure_account_key")

#### azure_account_key Type

`string`

## Definitions group smb_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters"}
```

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                         |
| :------------------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [smb_host](#smb_host)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")     |
| [smb_user](#smb_user)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")     |
| [smb_pass](#smb_pass)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")     |
| [smb_domain](#smb_domain) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain") |

### smb_host

Host name or IP address

`smb_host`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")

#### smb_host Type

`string`

#### smb_host Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")

### smb_user

User name for share connection

`smb_user`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")

#### smb_user Type

`string`

#### smb_user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### smb_pass

User password for share connection

`smb_pass`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")

#### smb_pass Type

`string`

### smb_domain

The short form (NT-style) domain name

`smb_domain`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain")

#### smb_domain Type

`string`

## Definitions group cluster_parameters

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/cluster_parameters"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |
