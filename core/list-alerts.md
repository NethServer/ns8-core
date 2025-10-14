# list-alerts output Schema

```txt
http://schema.nethserver.org/cluster/list-alerts.json
```

Output schema of the list-alerts action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-alerts.json](cluster/list-alerts.json "open original schema") |

## list-alerts output Type

`object` ([list-alerts output](list-alerts.md))

# list-alerts output Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                        |
| :---------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------- |
| [filter](#filter) | `object` | Optional | cannot be null | [list-alerts output](list-alerts-properties-filter.md "http://schema.nethserver.org/cluster/list-alerts.json#/properties/filter") |
| [alerts](#alerts) | `array`  | Required | cannot be null | [list-alerts output](list-alerts-properties-alerts.md "http://schema.nethserver.org/cluster/list-alerts.json#/properties/alerts") |

## filter



`filter`

* is optional

* Type: `object` ([Details](list-alerts-properties-filter.md))

* cannot be null

* defined in: [list-alerts output](list-alerts-properties-filter.md "http://schema.nethserver.org/cluster/list-alerts.json#/properties/filter")

### filter Type

`object` ([Details](list-alerts-properties-filter.md))

## alerts



`alerts`

* is required

* Type: an array of merged types ([Details](list-alerts-properties-alerts-items.md))

* cannot be null

* defined in: [list-alerts output](list-alerts-properties-alerts.md "http://schema.nethserver.org/cluster/list-alerts.json#/properties/alerts")

### alerts Type

an array of merged types ([Details](list-alerts-properties-alerts-items.md))

# list-alerts output Definitions

## Definitions group alertmanager-v2-object

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object"}
```

| Property                      | Type     | Required | Nullable       | Defined by                                                                                                                                                                                             |
| :---------------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [labels](#labels)             | `object` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-labels.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/labels")             |
| [updatedAt](#updatedat)       | `string` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-updatedat.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/updatedAt")       |
| [startsAt](#startsat)         | `string` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-startsat.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/startsAt")         |
| [endsAt](#endsat)             | `string` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-endsat.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/endsAt")             |
| [fingerprint](#fingerprint)   | `string` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-fingerprint.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/fingerprint")   |
| [status](#status)             | `object` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-status.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/status")             |
| [generatorURL](#generatorurl) | `string` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-generatorurl.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/generatorURL") |
| [annotations](#annotations)   | `object` | Optional | cannot be null | [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-annotations.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/annotations")   |

### labels



`labels`

* is optional

* Type: `object` ([Details](list-alerts-defs-alertmanager-v2-object-properties-labels.md))

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-labels.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/labels")

#### labels Type

`object` ([Details](list-alerts-defs-alertmanager-v2-object-properties-labels.md))

### updatedAt



`updatedAt`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-updatedat.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/updatedAt")

#### updatedAt Type

`string`

#### updatedAt Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

### startsAt



`startsAt`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-startsat.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/startsAt")

#### startsAt Type

`string`

#### startsAt Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

### endsAt



`endsAt`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-endsat.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/endsAt")

#### endsAt Type

`string`

#### endsAt Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

### fingerprint



`fingerprint`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-fingerprint.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/fingerprint")

#### fingerprint Type

`string`

### status



`status`

* is optional

* Type: `object` ([Details](list-alerts-defs-alertmanager-v2-object-properties-status.md))

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-status.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/status")

#### status Type

`object` ([Details](list-alerts-defs-alertmanager-v2-object-properties-status.md))

### generatorURL



`generatorURL`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-generatorurl.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/generatorURL")

#### generatorURL Type

`string`

### annotations



`annotations`

* is optional

* Type: `object` ([Details](list-alerts-defs-alertmanager-v2-object-properties-annotations.md))

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-alertmanager-v2-object-properties-annotations.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/annotations")

#### annotations Type

`object` ([Details](list-alerts-defs-alertmanager-v2-object-properties-annotations.md))

## Definitions group base-attrs

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs"}
```

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                     |
| :------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name) | `string` | Required | cannot be null | [list-alerts output](list-alerts-defs-base-attrs-properties-name.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/properties/name") |

### name

Short alert name

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-base-attrs-properties-name.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/properties/name")

#### name Type

`string`
