# retrieve-cluster-backup input Schema

```txt
http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json
```

Retrieve cluster backup from base64 field or download it from a HTTP URL

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [retrieve-cluster-backup-input.json](cluster/retrieve-cluster-backup-input.json "open original schema") |

## retrieve-cluster-backup input Type

`object` ([retrieve-cluster-backup input](retrieve-cluster-backup-input.md))

## retrieve-cluster-backup input Examples

```json
{
  "type": "file",
  "password": "mypassword",
  "file": "jA0ECQMC2c...94FL/8y2KV"
}
```

```json
{
  "type": "url",
  "password": "mypassword",
  "url": "https://myserver.org/dump.json.gz.gpg",
  "tls_verify": true
}
```

# retrieve-cluster-backup input Properties

| Property                   | Type      | Required | Nullable       | Defined by                                                                                                                                                                                            |
| :------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [password](#password)      | `string`  | Required | cannot be null | [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-gpg-passhprase.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/password")            |
| [type](#type)              | `string`  | Required | cannot be null | [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-select-backup-retrieve-method.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/type") |
| [url](#url)                | `string`  | Optional | cannot be null | [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-httpss-url.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/url")                     |
| [tls\_verify](#tls_verify) | `boolean` | Optional | cannot be null | [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-tls-verfication-flag.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/tls_verify")    |
| [file](#file)              | `string`  | Optional | cannot be null | [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-backup-file.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/file")                   |

## password

Passhprase for backup encryption

`password`

* is required

* Type: `string` ([GPG passhprase](retrieve-cluster-backup-input-properties-gpg-passhprase.md))

* cannot be null

* defined in: [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-gpg-passhprase.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/password")

### password Type

`string` ([GPG passhprase](retrieve-cluster-backup-input-properties-gpg-passhprase.md))

## type

Choose if the backup is uploaded inside the payload or should be downloaded from a URL. Can be 'file' or 'url'

`type`

* is required

* Type: `string` ([Select backup retrieve method](retrieve-cluster-backup-input-properties-select-backup-retrieve-method.md))

* cannot be null

* defined in: [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-select-backup-retrieve-method.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/type")

### type Type

`string` ([Select backup retrieve method](retrieve-cluster-backup-input-properties-select-backup-retrieve-method.md))

## url

HTTP/s URL of remote backup

`url`

* is optional

* Type: `string` ([HTTP/ss URL](retrieve-cluster-backup-input-properties-httpss-url.md))

* cannot be null

* defined in: [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-httpss-url.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/url")

### url Type

`string` ([HTTP/ss URL](retrieve-cluster-backup-input-properties-httpss-url.md))

### url Constraints

**URI**: the string must be a URI, according to [RFC 3986](https://tools.ietf.org/html/rfc3986 "check the specification")

## tls\_verify

Enable or disable TLS verification

`tls_verify`

* is optional

* Type: `boolean` ([TLS verfication flag](retrieve-cluster-backup-input-properties-tls-verfication-flag.md))

* cannot be null

* defined in: [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-tls-verfication-flag.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/tls_verify")

### tls\_verify Type

`boolean` ([TLS verfication flag](retrieve-cluster-backup-input-properties-tls-verfication-flag.md))

## file

Backup file encoded with base64

`file`

* is optional

* Type: `string` ([Backup file](retrieve-cluster-backup-input-properties-backup-file.md))

* cannot be null

* defined in: [retrieve-cluster-backup input](retrieve-cluster-backup-input-properties-backup-file.md "http://schema.nethserver.org/cluster/retrieve-cluster-backup-input.json#/properties/file")

### file Type

`string` ([Backup file](retrieve-cluster-backup-input-properties-backup-file.md))

### file Constraints

**encoding**: the string content must be using the base64 content encoding.

**media type**: the media type of the contents of this string is: `application/octet-stream`
