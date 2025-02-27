# Untitled object in list-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-modules-output.json\*](cluster/list-modules-output.json "open original schema") |

## items Type

`object` ([Details](list-modules-output-items-properties-versions-items.md))

# items Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                           |
| :------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [tag](#tag)         | `string`  | Required | cannot be null | [list-modules output](list-modules-output-items-properties-versions-items-properties-tag.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items/properties/tag")         |
| [testing](#testing) | `boolean` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-versions-items-properties-testing.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items/properties/testing") |
| [labels](#labels)   | `object`  | Required | cannot be null | [list-modules output](list-modules-output-items-properties-versions-items-properties-labels.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items/properties/labels")   |

## tag

The package version which must be valid semantic version, like '1.0.0

`tag`

* is required

* Type: `string`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-versions-items-properties-tag.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items/properties/tag")

### tag Type

`string`

## testing

Set to 'true' if this version is pre-release according to semantic versioning (<https://semver.org>)

`testing`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-versions-items-properties-testing.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items/properties/testing")

### testing Type

`boolean`

## labels

Image label, see official specification <https://github.com/opencontainers/image-spec/blob/master/annotations.md>

`labels`

* is required

* Type: `object` ([Details](list-modules-output-items-properties-versions-items-properties-labels.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-versions-items-properties-labels.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions/items/properties/labels")

### labels Type

`object` ([Details](list-modules-output-items-properties-versions-items-properties-labels.md))
