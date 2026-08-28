# add-repository input Schema

```txt
http://schema.nethserver.org/cluster/add-repository-input.json
```

Input schema of the add-repository action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-repository-input.json](cluster/add-repository-input.json "open original schema") |

## add-repository input Type

`object` ([add-repository input](add-repository-input.md))

## add-repository input Examples

```json
{
  "name": "repository1",
  "status": true,
  "url": "https://repository1.nethserver.org/"
}
```

# add-repository input Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                              |
| :------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [name](#name)       | `string`  | Required | cannot be null | [add-repository input](add-repository-input-properties-name.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/name")       |
| [url](#url)         | `string`  | Required | cannot be null | [add-repository input](add-repository-input-properties-url.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/url")         |
| [testing](#testing) | `boolean` | Optional | cannot be null | [add-repository input](add-repository-input-properties-testing.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/testing") |
| [status](#status)   | `boolean` | Required | cannot be null | [add-repository input](add-repository-input-properties-status.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/status")   |

## name

Unique repository name

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [add-repository input](add-repository-input-properties-name.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/name")

### name Type

`string`

## url

Base URL of the repository

`url`

* is required

* Type: `string`

* cannot be null

* defined in: [add-repository input](add-repository-input-properties-url.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/url")

### url Type

`string`

## testing

Use testing releases to install new instances and update existing ones.

`testing`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [add-repository input](add-repository-input-properties-testing.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/testing")

### testing Type

`boolean`

## status

Enable or disable the repository. When 'true', the repository is enabled.

`status`

* is required

* Type: `boolean`

* cannot be null

* defined in: [add-repository input](add-repository-input-properties-status.md "http://schema.nethserver.org/cluster/add-repository-input.json#/properties/status")

### status Type

`boolean`
