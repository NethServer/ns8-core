# update-core input Schema

```txt
http://schema.nethserver.org/cluster/update-core-input.json
```

Input schema of the update-core action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [update-core-input.json](cluster/update-core-input.json "open original schema") |

## update-core input Type

`object` ([update-core input](update-core-input-1.md))

## update-core input Examples

```json
{
  "core_url": "ghcr.io/nethserver/core:1.2.0",
  "nodes": [
    1,
    2,
    3,
    4
  ]
}
```

```json
{
  "nodes": [
    1
  ]
}
```

# update-core input Properties

| Property               | Type      | Required | Nullable       | Defined by                                                                                                                                         |
| :--------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------- |
| [core\_url](#core_url) | `string`  | Optional | cannot be null | [update-core input](update-core-input-1-properties-core_url.md "http://schema.nethserver.org/cluster/update-core-input.json#/properties/core_url") |
| [nodes](#nodes)        | `array`   | Required | cannot be null | [update-core input](update-core-input-1-properties-nodes.md "http://schema.nethserver.org/cluster/update-core-input.json#/properties/nodes")       |
| [force](#force)        | `boolean` | Optional | cannot be null | [update-core input](update-core-input-1-properties-force.md "http://schema.nethserver.org/cluster/update-core-input.json#/properties/force")       |

## core\_url

Core image URL to download and install. If empty, search core URL inside repository metadata.

`core_url`

* is optional

* Type: `string`

* cannot be null

* defined in: [update-core input](update-core-input-1-properties-core_url.md "http://schema.nethserver.org/cluster/update-core-input.json#/properties/core_url")

### core\_url Type

`string`

## nodes

Identifiers of nodes where the selected image is installed

`nodes`

* is required

* Type: `integer[]`

* cannot be null

* defined in: [update-core input](update-core-input-1-properties-nodes.md "http://schema.nethserver.org/cluster/update-core-input.json#/properties/nodes")

### nodes Type

`integer[]`

### nodes Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

## force

Force the core update even if the given `core_url` is already present in the local Podman storage

`force`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [update-core input](update-core-input-1-properties-force.md "http://schema.nethserver.org/cluster/update-core-input.json#/properties/force")

### force Type

`boolean`
