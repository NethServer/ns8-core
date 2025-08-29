# add-module input Schema

```txt
http://schema.nethserver.org/cluster/add-module-input.json
```

Input schema of the add-module action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-module-input.json](cluster/add-module-input.json "open original schema") |

## add-module input Type

`object` ([add-module input](add-module-input-1.md))

## add-module input Examples

```json
{
  "image": "traefik",
  "node": 1
}
```

```json
{
  "image": "ghcr.io/nethserver/traefik:latest",
  "node": 1
}
```

# add-module input Properties

| Property                              | Type      | Required | Nullable       | Defined by                                                                                                                                                               |
| :------------------------------------ | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [check\_idle\_time](#check_idle_time) | `integer` | Optional | cannot be null | [add-module input](add-module-input-1-properties-agent-liveness-check-limit.md "http://schema.nethserver.org/cluster/add-module-input.json#/properties/check_idle_time") |
| [image](#image)                       | `string`  | Required | cannot be null | [add-module input](add-module-input-1-properties-module-image-name.md "http://schema.nethserver.org/cluster/add-module-input.json#/properties/image")                    |
| [node](#node)                         | `integer` | Required | cannot be null | [add-module input](add-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/add-module-input.json#/properties/node")                               |

## check\_idle\_time

Change the default check limit value (default 8 seconds)

`check_idle_time`

* is optional

* Type: `integer` ([Agent liveness check limit](add-module-input-1-properties-agent-liveness-check-limit.md))

* cannot be null

* defined in: [add-module input](add-module-input-1-properties-agent-liveness-check-limit.md "http://schema.nethserver.org/cluster/add-module-input.json#/properties/check_idle_time")

### check\_idle\_time Type

`integer` ([Agent liveness check limit](add-module-input-1-properties-agent-liveness-check-limit.md))

## image

Name of the module image to install. If the image is not a URL, install the latest version of the module. Otherwise, just install the provided URL.

`image`

* is required

* Type: `string` ([Module image name](add-module-input-1-properties-module-image-name.md))

* cannot be null

* defined in: [add-module input](add-module-input-1-properties-module-image-name.md "http://schema.nethserver.org/cluster/add-module-input.json#/properties/image")

### image Type

`string` ([Module image name](add-module-input-1-properties-module-image-name.md))

## node

Node identifier where the module has to be installed

`node`

* is required

* Type: `integer` ([Node ID](add-module-input-1-properties-node-id.md))

* cannot be null

* defined in: [add-module input](add-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/add-module-input.json#/properties/node")

### node Type

`integer` ([Node ID](add-module-input-1-properties-node-id.md))

### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`
