# update-modules input Schema

```txt
http://schema.nethserver.org/cluster/update-modules-input.json
```

Input schema of the update-modules action. All matching modules are updated to their own valid latest version. The limits min-core and min-from are considered. If no option is given, all cluster app instances are updated.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [update-modules-input.json](cluster/update-modules-input.json "open original schema") |

## update-modules input Type

`object` ([update-modules input](update-modules-input.md))

## update-modules input Examples

```json
{
  "instances": [
    "mymodule2",
    "mymodule3"
  ]
}
```

```json
{
  "modules": [
    "mymodule"
  ]
}
```

# update-modules input Properties

| Property                | Type    | Required | Nullable       | Defined by                                                                                                                                                  |
| :---------------------- | :------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [instances](#instances) | `array` | Optional | cannot be null | [update-modules input](update-modules-input-properties-instances.md "http://schema.nethserver.org/cluster/update-modules-input.json#/properties/instances") |
| [modules](#modules)     | `array` | Optional | cannot be null | [update-modules input](update-modules-input-properties-modules.md "http://schema.nethserver.org/cluster/update-modules-input.json#/properties/modules")     |

## instances

Limit the update to this list of instances. Only matching instances are updated.

`instances`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [update-modules input](update-modules-input-properties-instances.md "http://schema.nethserver.org/cluster/update-modules-input.json#/properties/instances")

### instances Type

`string[]`

## modules

Limit the update to this list of modules. Only instances of this module are updated.

`modules`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [update-modules input](update-modules-input-properties-modules.md "http://schema.nethserver.org/cluster/update-modules-input.json#/properties/modules")

### modules Type

`string[]`
