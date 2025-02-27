# create-module input Schema

```txt
http://schema.nethserver.org/agent/create-module-input.json
```

Input schema of the basic create-module action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [create-module-input.json](agent/create-module-input.json "open original schema") |

## create-module input Type

`object` ([create-module input](create-module-input.md))

## create-module input Examples

```json
{
  "images": [
    "docker.io/traefik:v2.4"
  ]
}
```

# create-module input Properties

| Property          | Type    | Required | Nullable       | Defined by                                                                                                                                       |
| :---------------- | :------ | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------- |
| [images](#images) | `array` | Optional | cannot be null | [create-module input](create-module-input-properties-images.md "http://schema.nethserver.org/agent/create-module-input.json#/properties/images") |

## images



`images`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [create-module input](create-module-input-properties-images.md "http://schema.nethserver.org/agent/create-module-input.json#/properties/images")

### images Type

`string[]`

### images Constraints

**unique items**: all items in this array must be unique. Duplicates are not allowed.
