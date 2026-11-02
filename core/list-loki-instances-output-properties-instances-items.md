# Untitled object in List loki instances Schema

```txt
http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-loki-instances-output.json\*](cluster/list-loki-instances-output.json "open original schema") |

## items Type

`object` ([Details](list-loki-instances-output-properties-instances-items.md))

# items Properties

| Property                                | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                   |
| :-------------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [instance_id](#instance_id)             | `string`  | Required | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-instance_id.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/instance_id")             |
| [instance_label](#instance_label)       | `string`  | Required | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-instance_label.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/instance_label")       |
| [node_id](#node_id)                     | `string`  | Required | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-node_id.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/node_id")                     |
| [node_label](#node_label)               | `string`  | Required | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-node_label.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/node_label")               |
| [active](#active)                       | `boolean` | Required | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-active.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/active")                       |
| [offline](#offline)                     | `boolean` | Required | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-offline.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/offline")                     |
| [retention_days](#retention_days)       | `integer` | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-retention_days.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/retention_days")       |
| [active_from](#active_from)             | `string`  | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-active_from.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/active_from")             |
| [active_to](#active_to)                 | `string`  | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-active_to.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/active_to")                 |
| [cloud_log_manager](#cloud_log_manager) | `object`  | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager") |
| [syslog](#syslog)                       | `object`  | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-syslog.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/syslog")                       |

## instance_id

The Loki instance identifier.

`instance_id`

* is required

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-instance_id.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/instance_id")

### instance_id Type

`string`

## instance_label

The Loki instance label, if empty string, there's no label.

`instance_label`

* is required

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-instance_label.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/instance_label")

### instance_label Type

`string`

## node_id

The node identifier where the Loki instance is running.

`node_id`

* is required

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-node_id.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/node_id")

### node_id Type

`string`

## node_label

The node label where the Loki instance is running, can be empty.

`node_label`

* is required

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-node_label.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/node_label")

### node_label Type

`string`

## active

The Loki instance is the currently active.

`active`

* is required

* Type: `boolean`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-active.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/active")

### active Type

`boolean`

## offline

The Loki instance is offline.

`offline`

* is required

* Type: `boolean`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-offline.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/offline")

### offline Type

`boolean`

## retention_days

The retention days for the Loki instance.

`retention_days`

* is optional

* Type: `integer`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-retention_days.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/retention_days")

### retention_days Type

`integer`

## active_from

The ISO 8601 date-time when the Loki instance was activated.

`active_from`

* is optional

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-active_from.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/active_from")

### active_from Type

`string`

### active_from Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## active_to

The ISO 8601 date-time when the Loki instance was deactivated.

`active_to`

* is optional

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-active_to.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/active_to")

### active_to Type

`string`

### active_to Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")

## cloud_log_manager

Cloud Log Manager forwarder configuration.

`cloud_log_manager`

* is optional

* Type: `object` ([Details](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager.md))

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager")

### cloud_log_manager Type

`object` ([Details](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager.md))

## syslog

Syslog forwarder configuration.

`syslog`

* is optional

* Type: `object` ([Details](list-loki-instances-output-properties-instances-items-properties-syslog.md))

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-syslog.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/syslog")

### syslog Type

`object` ([Details](list-loki-instances-output-properties-instances-items-properties-syslog.md))
