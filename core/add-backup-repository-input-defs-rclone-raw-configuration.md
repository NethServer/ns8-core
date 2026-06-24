# Rclone raw configuration Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## rclone\_parameters Type

`object` ([Rclone raw configuration](add-backup-repository-input-defs-rclone-raw-configuration.md))

# rclone\_parameters Properties

| Property                                             | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                       |
| :--------------------------------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [rclone\_conf\_secret](#rclone_conf_secret)          | `string` | Required | can be null    | [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_secret")         |
| [rclone\_conf\_b64\_secret](#rclone_conf_b64_secret) | `string` | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_b64_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_b64_secret") |
| [basepath](#basepath)                                | `string` | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-basepath.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/basepath")                             |

## rclone\_conf\_secret

Basic INI-like string validation

`rclone_conf_secret`

* is required

* Type: `string`

* can be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_secret")

### rclone\_conf\_secret Type

`string`

### rclone\_conf\_secret Constraints

**maximum length**: the maximum number of characters for this string is: `100000`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^(\[[^\]]+\]\n)?(\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*[^\n]*(\n|$))*$
```

[try pattern](https://regexr.com/?expression=%5E\(%5C%5B%5B%5E%5C%5D%5D%2B%5C%5D%5Cn\)%3F\(%5Cs*\(%5BA-Za-z_%5D%5BA-Za-z0-9_%5D*\)%5Cs*%3D%5Cs*%5B%5E%5Cn%5D*\(%5Cn%7C%24\)\)*%24 "try regular expression with regexr.com")

## rclone\_conf\_b64\_secret



`rclone_conf_b64_secret`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-rclone_conf_b64_secret.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_b64_secret")

### rclone\_conf\_b64\_secret Type

`string`

### rclone\_conf\_b64\_secret Constraints

**maximum length**: the maximum number of characters for this string is: `100000`

## basepath



`basepath`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-rclone-raw-configuration-properties-basepath.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/$defs/rclone_parameters/properties/basepath")

### basepath Type

`string`

### basepath Constraints

**maximum length**: the maximum number of characters for this string is: `512`
