# Core image URL Schema

```txt
http://schema.nethserver.org/node/update-core-input.json#/properties/core_url
```

URL of the new core image

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                     |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [update-core-input.json\*](node/update-core-input.json "open original schema") |

## core\_url Type

`string` ([Core image URL](update-core-input-properties-core-image-url.md))

## core\_url Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## core\_url Examples

```json
"ghcr.io/nethserver/core:1.2.0"
```
