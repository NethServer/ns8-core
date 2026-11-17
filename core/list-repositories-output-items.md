# Untitled object in list-repositories output Schema

```txt
http://schema.nethserver.org/cluster/list-repositories-output.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-repositories-output.json\*](cluster/list-repositories-output.json "open original schema") |

## items Type

`object` ([Details](list-repositories-output-items.md))

# items Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                                                      |
| :------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [name](#name)       | `string`  | Required | cannot be null | [list-repositories output](list-repositories-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/name")       |
| [url](#url)         | `string`  | Required | cannot be null | [list-repositories output](list-repositories-output-items-properties-url.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/url")         |
| [testing](#testing) | `boolean` | Required | cannot be null | [list-repositories output](list-repositories-output-items-properties-testing.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/testing") |
| [status](#status)   | `boolean` | Required | cannot be null | [list-repositories output](list-repositories-output-items-properties-status.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/status")   |

## name

Unique repository name

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-repositories output](list-repositories-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/name")

### name Type

`string`

## url

Base URL of the repository

`url`

* is required

* Type: `string`

* cannot be null

* defined in: [list-repositories output](list-repositories-output-items-properties-url.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/url")

### url Type

`string`

## testing

Enable or disable access to testing images. When 'true', the repository will list testing images.

`testing`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-repositories output](list-repositories-output-items-properties-testing.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/testing")

### testing Type

`boolean`

## status

Enable or disable the repository. When 'true', the repository is enabled.

`status`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-repositories output](list-repositories-output-items-properties-status.md "http://schema.nethserver.org/cluster/list-repositories-output.json#/items/properties/status")

### status Type

`boolean`
