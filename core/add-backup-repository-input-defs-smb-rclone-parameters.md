# SMB Rclone parameters Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## smb_parameters Type

`object` ([SMB Rclone parameters](add-backup-repository-input-defs-smb-rclone-parameters.md))

# smb_parameters Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                         |
| :------------------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [smb_host](#smb_host)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")     |
| [smb_user](#smb_user)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")     |
| [smb_pass](#smb_pass)     | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")     |
| [smb_domain](#smb_domain) | `string` | Required | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain") |

## smb_host

Host name or IP address

`smb_host`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")

### smb_host Type

`string`

### smb_host Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")

## smb_user

User name for share connection

`smb_user`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")

### smb_user Type

`string`

### smb_user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## smb_pass

User password for share connection

`smb_pass`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")

### smb_pass Type

`string`

## smb_domain

The short form (NT-style) domain name

`smb_domain`

* is required

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain")

### smb_domain Type

`string`
