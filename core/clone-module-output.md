# clone-module output Schema

```txt
http://schema.nethserver.org/cluster/clone-module-output.json
```

Output schema of the clone-module action. It has the same format of the add-module action output

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [clone-module-output.json](cluster/clone-module-output.json "open original schema") |

## clone-module output Type

`object` ([clone-module output](clone-module-output.md))

## clone-module output Examples

```json
{
  "module_id": "traefik1",
  "image_name": "Traefik edge proxy",
  "image_url": "ghcr.io/nethserver/traefik:latest"
}
```

# clone-module output Properties

| Property                    | Type     | Required | Nullable       | Defined by                                                                                                                                                   |
| :-------------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module_id](#module_id)     | `string` | Required | cannot be null | [clone-module output](clone-module-output-properties-module_id.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/module_id")     |
| [module_uuid](#module_uuid) | `string` | Optional | cannot be null | [clone-module output](clone-module-output-properties-module_uuid.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/module_uuid") |
| [image_name](#image_name)   | `string` | Required | cannot be null | [clone-module output](clone-module-output-properties-image_name.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/image_name")   |
| [image_url](#image_url)     | `string` | Required | cannot be null | [clone-module output](clone-module-output-properties-image_url.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/image_url")     |

## module_id

Generated identifier of the clone instance

`module_id`

* is required

* Type: `string`

* cannot be null

* defined in: [clone-module output](clone-module-output-properties-module_id.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/module_id")

### module_id Type

`string`

## module_uuid

Generated UUID of the added module instance

`module_uuid`

* is optional

* Type: `string`

* cannot be null

* defined in: [clone-module output](clone-module-output-properties-module_uuid.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/module_uuid")

### module_uuid Type

`string`

## image_name

Display name of the clone image

`image_name`

* is required

* Type: `string`

* cannot be null

* defined in: [clone-module output](clone-module-output-properties-image_name.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/image_name")

### image_name Type

`string`

## image_url

Repository URL of the clone image

`image_url`

* is required

* Type: `string`

* cannot be null

* defined in: [clone-module output](clone-module-output-properties-image_url.md "http://schema.nethserver.org/cluster/clone-module-output.json#/properties/image_url")

### image_url Type

`string`
