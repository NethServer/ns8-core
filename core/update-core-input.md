# update-core input Schema

```txt
http://schema.nethserver.org/node/update-core-input.json
```

Update the core module on the local node

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                   |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [update-core-input.json](node/update-core-input.json "open original schema") |

## update-core input Type

`object` ([update-core input](update-core-input.md))

## update-core input Examples

```json
{
  "core_url": "ghcr.io/nethserver/core:1.2.0"
}
```

# update-core input Properties

| Property               | Type      | Required | Nullable       | Defined by                                                                                                                                          |
| :--------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| [core\_url](#core_url) | `string`  | Required | cannot be null | [update-core input](update-core-input-properties-core-image-url.md "http://schema.nethserver.org/node/update-core-input.json#/properties/core_url") |
| [force](#force)        | `boolean` | Optional | cannot be null | [update-core input](update-core-input-properties-force.md "http://schema.nethserver.org/node/update-core-input.json#/properties/force")             |

## core\_url

URL of the new core image

`core_url`

* is required

* Type: `string` ([Core image URL](update-core-input-properties-core-image-url.md))

* cannot be null

* defined in: [update-core input](update-core-input-properties-core-image-url.md "http://schema.nethserver.org/node/update-core-input.json#/properties/core_url")

### core\_url Type

`string` ([Core image URL](update-core-input-properties-core-image-url.md))

### core\_url Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### core\_url Examples

```json
"ghcr.io/nethserver/core:1.2.0"
```

## force

Force the core update even if the given `core_url` is already present in the local Podman storage

`force`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [update-core input](update-core-input-properties-force.md "http://schema.nethserver.org/node/update-core-input.json#/properties/force")

### force Type

`boolean`
