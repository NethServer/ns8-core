# update-module input Schema

```txt
http://schema.nethserver.org/agent/update-module-input.json
```

Input schema of the basic update-module action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [update-module-input.json](agent/update-module-input.json "open original schema") |

## update-module input Type

`object` ([update-module input](update-module-input-1.md))

## update-module input Examples

```json
{
  "module_url": "ghcr.io/nethserver/mymodule:3.2.1"
}
```

# update-module input Properties

| Property                   | Type      | Required | Nullable       | Defined by                                                                                                                                                 |
| :------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_url](#module_url) | `string`  | Optional | cannot be null | [update-module input](update-module-input-1-properties-module_url.md "http://schema.nethserver.org/agent/update-module-input.json#/properties/module_url") |
| [force](#force)            | `boolean` | Optional | cannot be null | [update-module input](update-module-input-1-properties-force.md "http://schema.nethserver.org/agent/update-module-input.json#/properties/force")           |

## module\_url



`module_url`

* is optional

* Type: `string`

* cannot be null

* defined in: [update-module input](update-module-input-1-properties-module_url.md "http://schema.nethserver.org/agent/update-module-input.json#/properties/module_url")

### module\_url Type

`string`

### module\_url Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## force



`force`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [update-module input](update-module-input-1-properties-force.md "http://schema.nethserver.org/agent/update-module-input.json#/properties/force")

### force Type

`boolean`
