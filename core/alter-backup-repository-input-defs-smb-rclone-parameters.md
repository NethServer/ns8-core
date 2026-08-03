# SMB Rclone parameters Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## smb\_parameters Type

`object` ([SMB Rclone parameters](alter-backup-repository-input-defs-smb-rclone-parameters.md))

# smb\_parameters Properties

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                               |
| :------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [smb\_host](#smb_host)     | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")     |
| [smb\_user](#smb_user)     | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")     |
| [smb\_pass](#smb_pass)     | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")     |
| [smb\_domain](#smb_domain) | `string` | Required | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain") |

## smb\_host

Host name or IP address

`smb_host`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_host.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host")

### smb\_host Type

`string`

### smb\_host Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")

## smb\_user

User name for share connection

`smb_user`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_user.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_user")

### smb\_user Type

`string`

### smb\_user Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## smb\_pass

User password for share connection

`smb_pass`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_pass.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_pass")

### smb\_pass Type

`string`

## smb\_domain

The short form (NT-style) domain name

`smb_domain`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-smb-rclone-parameters-properties-smb_domain.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_domain")

### smb\_domain Type

`string`
