# restore-module input Schema

```txt
http://schema.nethserver.org/module/restore-module-input.json
```

Restore the module state from a remote backup

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                           |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [restore-module-input.json](module/restore-module-input.json "open original schema") |

## restore-module input Type

`object` ([restore-module input](restore-module-input.md))

## restore-module input Examples

```json
{
  "repository": "48ce000a-79b7-5fe6-8558-177fd70c27b4",
  "path": "dokuwiki/dokuwiki1@f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf",
  "snapshot": "",
  "environment": {
    "IMAGE_URL": "ghcr.io/nethserver/dokuwiki:latest",
    "MODULE_UUID": "f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf",
    "DOKUWIKI_IMAGE": "docker.io/bitnami/dokuwiki:20200729.0.0-debian-10-r299",
    "DOKUWIKI_WIKI_NAME": "mywiki",
    "DOKUWIKI_USERNAME": "admin",
    "DOKUWIKI_PASSWORD": "pass",
    "DOKUWIKI_EMAIL": "davidep@nethesis.it",
    "DOKUWIKI_FULL_NAME": "Wiki Admin",
    "PHP_ENABLE_OPCACHE": "1",
    "PHP_MEMORY_LIMIT": "512M",
    "TRAEFIK_HOST": "mywiki.dp.nethserver.net",
    "TRAEFIK_HTTP2HTTPS": "False",
    "TRAEFIK_LETS_ENCRYPT": "False"
  },
  "replace": false
}
```

# restore-module input Properties

| Property                    | Type      | Required | Nullable       | Defined by                                                                                                                                                            |
| :-------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [environment](#environment) | `object`  | Required | cannot be null | [restore-module input](restore-module-input-properties-environment-backup.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/environment") |
| [repository](#repository)   | `string`  | Required | cannot be null | [restore-module input](restore-module-input-properties-repository-id.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/repository")       |
| [path](#path)               | `string`  | Required | cannot be null | [restore-module input](restore-module-input-properties-backup-path.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/path")               |
| [snapshot](#snapshot)       | `string`  | Required | cannot be null | [restore-module input](restore-module-input-properties-snapshot-id.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/snapshot")           |
| [replace](#replace)         | `boolean` | Optional | cannot be null | [restore-module input](restore-module-input-properties-replace-flag.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/replace")           |

## environment

Environment restored from the given backup

`environment`

* is required

* Type: `object` ([Environment backup](restore-module-input-properties-environment-backup.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-properties-environment-backup.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/environment")

### environment Type

`object` ([Environment backup](restore-module-input-properties-environment-backup.md))

## repository

Backup source repository

`repository`

* is required

* Type: `string` ([Repository ID](restore-module-input-properties-repository-id.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-properties-repository-id.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/repository")

### repository Type

`string` ([Repository ID](restore-module-input-properties-repository-id.md))

### repository Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## path

Path of the backup in the source repository

`path`

* is required

* Type: `string` ([Backup path](restore-module-input-properties-backup-path.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-properties-backup-path.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/path")

### path Type

`string` ([Backup path](restore-module-input-properties-backup-path.md))

### path Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## snapshot



`snapshot`

* is required

* Type: `string` ([Snapshot ID](restore-module-input-properties-snapshot-id.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-properties-snapshot-id.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/snapshot")

### snapshot Type

`string` ([Snapshot ID](restore-module-input-properties-snapshot-id.md))

## replace

If set to true restore old UUID from 'environment' field

`replace`

* is optional

* Type: `boolean` ([Replace flag](restore-module-input-properties-replace-flag.md))

* cannot be null

* defined in: [restore-module input](restore-module-input-properties-replace-flag.md "http://schema.nethserver.org/module/restore-module-input.json#/properties/replace")

### replace Type

`boolean` ([Replace flag](restore-module-input-properties-replace-flag.md))
