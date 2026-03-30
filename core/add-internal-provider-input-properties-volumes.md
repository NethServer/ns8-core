# Untitled object in add-internal-provider input Schema

```txt
http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/volumes
```

Optional mapping of application volume names to filesystem mount points.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-internal-provider-input.json\*](cluster/add-internal-provider-input.json "open original schema") |

## volumes Type

`object` ([Details](add-internal-provider-input-properties-volumes.md))

# volumes Properties

| Property | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                               |
| :------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `[^/]+`  | `string` | Optional | cannot be null | [add-internal-provider input](add-internal-provider-input-properties-volumes-patternproperties-.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/volumes/patternProperties/\[^/]+") |

## Pattern: `[^/]+`

Absolute filesystem directory path under /

`[^/]+`

* is optional

* Type: `string`

* cannot be null

* defined in: [add-internal-provider input](add-internal-provider-input-properties-volumes-patternproperties-.md "http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/volumes/patternProperties/\[^/]+")

### ]+ Type

`string`

### ]+ Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")
