# add-module output Schema

```txt
http://schema.nethserver.org/cluster/add-module-output.json
```

Output schema of the add-module action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-module-output.json](cluster/add-module-output.json "open original schema") |

## add-module output Type

`object` ([add-module output](add-module-output-1.md))

## add-module output Examples

```json
{
  "module_id": "traefik1",
  "image_name": "Traefik edge proxy",
  "image_url": "ghcr.io/nethserver/traefik:latest"
}
```

# add-module output Properties

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                             |
| :------------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id)   | `string` | Required | cannot be null | [add-module output](add-module-output-1-properties-module_id.md "http://schema.nethserver.org/cluster/add-module-output.json#/properties/module_id")   |
| [image\_name](#image_name) | `string` | Required | cannot be null | [add-module output](add-module-output-1-properties-image_name.md "http://schema.nethserver.org/cluster/add-module-output.json#/properties/image_name") |
| [image\_url](#image_url)   | `string` | Required | cannot be null | [add-module output](add-module-output-1-properties-image_url.md "http://schema.nethserver.org/cluster/add-module-output.json#/properties/image_url")   |

## module\_id

Generated identifier of the added module instance

`module_id`

* is required

* Type: `string`

* cannot be null

* defined in: [add-module output](add-module-output-1-properties-module_id.md "http://schema.nethserver.org/cluster/add-module-output.json#/properties/module_id")

### module\_id Type

`string`

## image\_name

Display name of the installed image

`image_name`

* is required

* Type: `string`

* cannot be null

* defined in: [add-module output](add-module-output-1-properties-image_name.md "http://schema.nethserver.org/cluster/add-module-output.json#/properties/image_name")

### image\_name Type

`string`

## image\_url

Repository URL of the installed image

`image_url`

* is required

* Type: `string`

* cannot be null

* defined in: [add-module output](add-module-output-1-properties-image_url.md "http://schema.nethserver.org/cluster/add-module-output.json#/properties/image_url")

### image\_url Type

`string`
