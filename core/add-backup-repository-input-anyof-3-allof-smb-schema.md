# SMB schema Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/3/allOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## 0 Type

unknown ([SMB schema](add-backup-repository-input-anyof-3-allof-smb-schema.md))

# 0 Properties

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                  |
| :------------------------ | :------------ | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [provider](#provider)     | Not specified | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-anyof-3-allof-smb-schema-properties-smb-provider.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/3/allOf/0/properties/provider") |
| [url](#url)               | Not specified | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-anyof-3-allof-smb-schema-properties-url.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/3/allOf/0/properties/url")               |
| [parameters](#parameters) | `object`      | Optional | cannot be null | [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/3/allOf/0/properties/parameters")                     |

## provider



`provider`

* is optional

* Type: unknown ([SMB provider](add-backup-repository-input-anyof-3-allof-smb-schema-properties-smb-provider.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-anyof-3-allof-smb-schema-properties-smb-provider.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/3/allOf/0/properties/provider")

### provider Type

unknown ([SMB provider](add-backup-repository-input-anyof-3-allof-smb-schema-properties-smb-provider.md))

### provider Constraints

**constant**: the value of this property must be equal to:

```json
"smb"
```

## url



`url`

* is optional

* Type: unknown

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-anyof-3-allof-smb-schema-properties-url.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/3/allOf/0/properties/url")

### url Type

unknown

### url Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^smb:.+$
```

[try pattern](https://regexr.com/?expression=%5Esmb%3A.%2B%24 "try regular expression with regexr.com")

## parameters



`parameters`

* is optional

* Type: `object` ([SMB Rclone parameters](add-backup-repository-input-defs-smb-rclone-parameters.md))

* cannot be null

* defined in: [add-backup-repository input](add-backup-repository-input-defs-smb-rclone-parameters.md "http://schema.nethserver.org/cluster/add-backup-repository-input.json#/anyOf/3/allOf/0/properties/parameters")

### parameters Type

`object` ([SMB Rclone parameters](add-backup-repository-input-defs-smb-rclone-parameters.md))
