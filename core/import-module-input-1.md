# import-module input Schema

```txt
http://schema.nethserver.org/cluster/import-module-input.json
```

Input schema of the import-module action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [import-module-input.json](cluster/import-module-input.json "open original schema") |

## import-module input Type

`object` ([import-module input](import-module-input-1.md))

## import-module input Examples

```json
{
  "image": "traefik",
  "node": 1,
  "volumes": [
    "traefik-data"
  ]
}
```

```json
{
  "image": "ghcr.io/nethserver/traefik:latest",
  "node": 1,
  "volumes": [
    "traefik-data"
  ]
}
```

# import-module input Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                                                                         |
| :------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [volumes](#volumes) | `array`   | Required | cannot be null | [import-module input](import-module-input-1-properties-initial-volume-set-where-the-module-state-is-stored.md "http://schema.nethserver.org/cluster/import-module-input.json#/properties/volumes") |
| [image](#image)     | `string`  | Required | cannot be null | [import-module input](import-module-input-1-properties-module-image-name.md "http://schema.nethserver.org/cluster/import-module-input.json#/properties/image")                                     |
| [node](#node)       | `integer` | Required | cannot be null | [import-module input](import-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/import-module-input.json#/properties/node")                                                |

## volumes



`volumes`

* is required

* Type: `string[]` ([Name of the volume element](import-module-input-1-properties-initial-volume-set-where-the-module-state-is-stored-name-of-the-volume-element.md))

* cannot be null

* defined in: [import-module input](import-module-input-1-properties-initial-volume-set-where-the-module-state-is-stored.md "http://schema.nethserver.org/cluster/import-module-input.json#/properties/volumes")

### volumes Type

`string[]` ([Name of the volume element](import-module-input-1-properties-initial-volume-set-where-the-module-state-is-stored-name-of-the-volume-element.md))

## image

Name of the module image to install. If the image is not a URL, install the latest version of the module. Otherwise, just install the provided URL.

`image`

* is required

* Type: `string` ([Module image name](import-module-input-1-properties-module-image-name.md))

* cannot be null

* defined in: [import-module input](import-module-input-1-properties-module-image-name.md "http://schema.nethserver.org/cluster/import-module-input.json#/properties/image")

### image Type

`string` ([Module image name](import-module-input-1-properties-module-image-name.md))

## node

Node identifier where the module has to be installed

`node`

* is required

* Type: `integer` ([Node ID](import-module-input-1-properties-node-id.md))

* cannot be null

* defined in: [import-module input](import-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/import-module-input.json#/properties/node")

### node Type

`integer` ([Node ID](import-module-input-1-properties-node-id.md))

### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`
