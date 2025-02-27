# Unconfigured domain Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain
```

An account provider instance, installed as the first instance of a new domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-user-domains-output.json\*](cluster/list-user-domains-output.json "open original schema") |

## unconfigured-domain Type

`object` ([Unconfigured domain](list-user-domains-output-defs-unconfigured-domain.md))

# unconfigured-domain Properties

| Property                 | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                              |
| :----------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id) | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/module_id")      |
| [node](#node)            | `integer` | Required | can be null    | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/node")             |
| [location](#location)    | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/location") |
| [protocol](#protocol)    | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/protocol")       |
| [schema](#schema)        | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/schema")      |

## module\_id

e.g. `samba1`

`module_id`

* is required

* Type: `string` ([Module identifier](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/module_id")

### module\_id Type

`string` ([Module identifier](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md))

## node

The node number, e.g. `1`

`node`

* is required

* Type: `integer` ([Node identifier](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md))

* can be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/node")

### node Type

`integer` ([Node identifier](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md))

### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## location

Type of domain hosting. Set to `internal` if the domain is hosted by the cluster, `external` if the domain is provided by a remote service

`location`

* is required

* Type: `string` ([Domain hosting location](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/location")

### location Type

`string` ([Domain hosting location](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md))

### location Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value        | Explanation |
| :----------- | :---------- |
| `"internal"` |             |
| `"external"` |             |

## protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/protocol")

### protocol Type

`string` ([Provider protocol](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

## schema



`schema`

* is required

* Type: `string` ([LDAP database schema](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/schema")

### schema Type

`string` ([LDAP database schema](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md))

### schema Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value       | Explanation |
| :---------- | :---------- |
| `"ad"`      |             |
| `"rfc2307"` |             |
