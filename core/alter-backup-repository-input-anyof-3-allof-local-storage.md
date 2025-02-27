# Local storage Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/3/allOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## 0 Type

unknown ([Local storage](alter-backup-repository-input-anyof-3-allof-local-storage.md))

# 0 Properties

| Property                  | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                                                        |
| :------------------------ | :------------ | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [provider](#provider)     | Not specified | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-anyof-3-allof-local-storage-properties-cluster-internal-provider.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/3/allOf/0/properties/provider") |
| [url](#url)               | Not specified | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-anyof-3-allof-local-storage-properties-url.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/3/allOf/0/properties/url")                            |
| [parameters](#parameters) | `object`      | Optional | cannot be null | [alter-backup-repository input](alter-backup-repository-input-defs-cluster-internal-rclone-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/3/allOf/0/properties/parameters")                        |

## provider



`provider`

* is optional

* Type: unknown ([Cluster internal provider](alter-backup-repository-input-anyof-3-allof-local-storage-properties-cluster-internal-provider.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-anyof-3-allof-local-storage-properties-cluster-internal-provider.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/3/allOf/0/properties/provider")

### provider Type

unknown ([Cluster internal provider](alter-backup-repository-input-anyof-3-allof-local-storage-properties-cluster-internal-provider.md))

### provider Constraints

**constant**: the value of this property must be equal to:

```json
"cluster"
```

## url



`url`

* is optional

* Type: unknown

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-anyof-3-allof-local-storage-properties-url.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/3/allOf/0/properties/url")

### url Type

unknown

### url Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^webdav:.+$
```

[try pattern](https://regexr.com/?expression=%5Ewebdav%3A.%2B%24 "try regular expression with regexr.com")

## parameters



`parameters`

* is optional

* Type: `object` ([Cluster-internal Rclone parameters](alter-backup-repository-input-defs-cluster-internal-rclone-parameters.md))

* cannot be null

* defined in: [alter-backup-repository input](alter-backup-repository-input-defs-cluster-internal-rclone-parameters.md "http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/3/allOf/0/properties/parameters")

### parameters Type

`object` ([Cluster-internal Rclone parameters](alter-backup-repository-input-defs-cluster-internal-rclone-parameters.md))
