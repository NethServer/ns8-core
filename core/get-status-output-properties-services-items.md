# Untitled object in Get module status Schema

```txt
http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-status-output.json\*](agent/get-status-output.json "open original schema") |

## items Type

`object` ([Details](get-status-output-properties-services-items.md))

# items Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                                                                       |
| :------------------ | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)       | `string`  | Optional | cannot be null | [Get module status](get-status-output-properties-services-items-properties-name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/name")       |
| [active](#active)   | `boolean` | Optional | cannot be null | [Get module status](get-status-output-properties-services-items-properties-active.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/active")   |
| [enabled](#enabled) | `boolean` | Optional | cannot be null | [Get module status](get-status-output-properties-services-items-properties-enabled.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/enabled") |
| [failed](#failed)   | `boolean` | Optional | cannot be null | [Get module status](get-status-output-properties-services-items-properties-failed.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/failed")   |

## name

Name of systemd unit

`name`

* is optional

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-services-items-properties-name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/name")

### name Type

`string`

## active

true if unit is running, false otherwise

`active`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [Get module status](get-status-output-properties-services-items-properties-active.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/active")

### active Type

`boolean`

## enabled

true if unit is enabled at startup, false otherwise

`enabled`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [Get module status](get-status-output-properties-services-items-properties-enabled.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/enabled")

### enabled Type

`boolean`

## failed

true if unit is in failed state, false otherwise

`failed`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [Get module status](get-status-output-properties-services-items-properties-failed.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services/items/properties/failed")

### failed Type

`boolean`
