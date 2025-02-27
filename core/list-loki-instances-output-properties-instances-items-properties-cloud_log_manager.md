# Untitled object in List loki instances Schema

```txt
http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager
```

Cloud Log Manager forwarder configuration.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-loki-instances-output.json\*](cluster/list-loki-instances-output.json "open original schema") |

## cloud\_log\_manager Type

`object` ([Details](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager.md))

# cloud\_log\_manager Properties

| Property                           | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                                       |
| :--------------------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [status](#status)                  | `string` | Required | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-status.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/status")                 |
| [address](#address)                | `string` | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-address.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/address")               |
| [tenant](#tenant)                  | `string` | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-tenant.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/tenant")                 |
| [last\_timestamp](#last_timestamp) | `string` | Optional | cannot be null | [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-last_timestamp.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/last_timestamp") |

## status

Forwarder state.

`status`

* is required

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-status.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/status")

### status Type

`string`

## address

Forwarder address where datas are sent.

`address`

* is optional

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-address.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/address")

### address Type

`string`

## tenant

Cloud Log Manager internal data.

`tenant`

* is optional

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-tenant.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/tenant")

### tenant Type

`string`

## last\_timestamp

Timestamp of the last log sent.

`last_timestamp`

* is optional

* Type: `string`

* cannot be null

* defined in: [List loki instances](list-loki-instances-output-properties-instances-items-properties-cloud_log_manager-properties-last_timestamp.md "http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/cloud_log_manager/properties/last_timestamp")

### last\_timestamp Type

`string`
