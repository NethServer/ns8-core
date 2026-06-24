# Untitled object in list-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason
```

If the versions array is empty, this object is present and explains why.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-modules-output.json\*](cluster/list-modules-output.json "open original schema") |

## no\_version\_reason Type

`object` ([Details](list-modules-output-items-properties-no_version_reason.md))

# no\_version\_reason Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                 |
| :------------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [message](#message) | `string` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-no_version_reason-properties-message.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason/properties/message") |
| [params](#params)   | `object` | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-no_version_reason-properties-params.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason/properties/params")   |

## message

Reason identifier

`message`

* is required

* Type: `string`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-no_version_reason-properties-message.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason/properties/message")

### message Type

`string`

## params

Parameters for the reason explanation

`params`

* is optional

* Type: `object` ([Details](list-modules-output-items-properties-no_version_reason-properties-params.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-no_version_reason-properties-params.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason/properties/params")

### params Type

`object` ([Details](list-modules-output-items-properties-no_version_reason-properties-params.md))
