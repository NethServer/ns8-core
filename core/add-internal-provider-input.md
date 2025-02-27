# add-internal-provider input Schema

```txt
http://schema.nethserver.org/cluster/add-internal-provider-input.json
```

Add a provider instance for a new or already existing internal user domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-internal-provider-input.json](cluster/add-internal-provider-input.json "open original schema") |

## add-internal-provider input Type

`object` ([add-internal-provider input](add-internal-provider-input.md))

## add-internal-provider input Examples

```json
{
  "image": "ghcr.io/nethserver/samba:v1.0.0",
  "node": 2
}
```

```json
{
  "image": "ghcr.io/nethserver/samba:v1.0.0",
  "node": 1,
  "domain": "ad.example.com"
}
```

# add-internal-provider input Properties

| Property          | Type      | Required | Nullable       | Defined by                                                                                                                                                                          |
| :---------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [image](#image)   | `string`  | Required | cannot be null | [add-internal-provider input](add-internal-provider-input-properties-module-image-url.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/image") |
| [node](#node)     | `integer` | Required | cannot be null | [add-internal-provider input](add-internal-provider-input-properties-node-id.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/node")           |
| [domain](#domain) | `string`  | Optional | can be null    | [add-internal-provider input](add-internal-provider-input-properties-domain.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/domain")          |

## image

URL of the module image to install

`image`

* is required

* Type: `string` ([Module image URL](add-internal-provider-input-properties-module-image-url.md))

* cannot be null

* defined in: [add-internal-provider input](add-internal-provider-input-properties-module-image-url.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/image")

### image Type

`string` ([Module image URL](add-internal-provider-input-properties-module-image-url.md))

### image Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## node

Node identifier where the module has to be installed

`node`

* is required

* Type: `integer` ([Node ID](add-internal-provider-input-properties-node-id.md))

* cannot be null

* defined in: [add-internal-provider input](add-internal-provider-input-properties-node-id.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/node")

### node Type

`integer` ([Node ID](add-internal-provider-input-properties-node-id.md))

### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## domain

Add the instance to an already existing domain. This field is optional.

`domain`

* is optional

* Type: `string` ([Domain](add-internal-provider-input-properties-domain.md))

* can be null

* defined in: [add-internal-provider input](add-internal-provider-input-properties-domain.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/domain")

### domain Type

`string` ([Domain](add-internal-provider-input-properties-domain.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

# add-internal-provider input Definitions
