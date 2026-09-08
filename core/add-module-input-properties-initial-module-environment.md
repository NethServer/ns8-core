# Initial module environment Schema

```txt
http://schema.nethserver.org/node/add-module-input.json#/properties/environment
```

Assign initial values to the module environment

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                   |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [add-module-input.json\*](node/add-module-input.json "open original schema") |

## environment Type

`object` ([Initial module environment](add-module-input-properties-initial-module-environment.md))

# environment Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                    |
| :----------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [IMAGE\_URL](#image_url) | `string` | Optional | cannot be null | [add-module input](add-module-input-properties-initial-module-environment-properties-image_url.md "http://schema.nethserver.org/node/add-module-input.json#/properties/environment/properties/IMAGE_URL")     |
| `^[^=]+$`                | `string` | Optional | cannot be null | [add-module input](add-module-input-properties-initial-module-environment-patternproperties-.md "http://schema.nethserver.org/node/add-module-input.json#/properties/environment/patternProperties/^\[^=]+$") |

## IMAGE\_URL

URL of the module root image

`IMAGE_URL`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-module input](add-module-input-properties-initial-module-environment-properties-image_url.md "http://schema.nethserver.org/node/add-module-input.json#/properties/environment/properties/IMAGE_URL")

### IMAGE\_URL Type

`string`

### IMAGE\_URL Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### IMAGE\_URL Examples

```json
"ghcr.io/nethserver/mymodule:v2.3.2"
```

## Pattern: `^[^=]+$`



`^[^=]+$`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-module input](add-module-input-properties-initial-module-environment-patternproperties-.md "http://schema.nethserver.org/node/add-module-input.json#/properties/environment/patternProperties/^\[^=]+$")

### ^\[^=]+$ Type

`string`
