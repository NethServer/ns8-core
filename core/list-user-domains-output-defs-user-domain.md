# User domain Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain
```

Users (and also user groups) can be uniquely identified inside a domain. An application is ususally bound to just one user domain at a time, but mulitple domains can be configured in the same cluster.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-user-domains-output.json\*](cluster/list-user-domains-output.json "open original schema") |

## user-domain Type

`object` ([User domain](list-user-domains-output-defs-user-domain.md))

any of

* not

  * [Protocol property is ldap](list-user-domains-output-defs-user-domain-anyof-0-protocol-property-is-ldap.md "check type definition")

* [LDAP domain properties](list-user-domains-output-defs-ldap-domain-properties.md "check type definition")

# user-domain Properties

| Property                         | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                              |
| :------------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)                    | `string` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-name.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/name")                        |
| [location](#location)            | `string` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/location") |
| [counters](#counters)            | `object` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters")                |
| [hidden\_users](#hidden_users)   | `array`  | Optional | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_users.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_users")        |
| [hidden\_groups](#hidden_groups) | `array`  | Optional | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_groups.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_groups")      |
| [protocol](#protocol)            | `string` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/protocol")       |
| [providers](#providers)          | `array`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-account-providers.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/providers")      |

## name



`name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-name.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/name")

### name Type

`string`

### name Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## location

Type of domain hosting. Set to `internal` if the domain is hosted by the cluster, `external` if the domain is provided by a remote service

`location`

* is required

* Type: `string` ([Domain hosting location](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/location")

### location Type

`string` ([Domain hosting location](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md))

### location Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value        | Explanation |
| :----------- | :---------- |
| `"internal"` |             |
| `"external"` |             |

## counters

The cached number of users and groups returned by their respective last API calls

`counters`

* is required

* Type: `object` ([Counters](list-user-domains-output-defs-user-domain-properties-counters.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters")

### counters Type

`object` ([Counters](list-user-domains-output-defs-user-domain-properties-counters.md))

## hidden\_users

List of users that are not visible from UI and from applications

`hidden_users`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_users.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_users")

### hidden\_users Type

`string[]`

## hidden\_groups

List of groups that are not visible from UI and from applications

`hidden_groups`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_groups.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_groups")

### hidden\_groups Type

`string[]`

## protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](list-user-domains-output-defs-user-domain-properties-provider-protocol.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/protocol")

### protocol Type

`string` ([Provider protocol](list-user-domains-output-defs-user-domain-properties-provider-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

## providers

Backend system and replicas providing the services of the user domain

`providers`

* is required

* Type: an array of merged types ([Details](list-user-domains-output-defs-user-domain-properties-account-providers-items.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-account-providers.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/providers")

### providers Type

an array of merged types ([Details](list-user-domains-output-defs-user-domain-properties-account-providers-items.md))

### providers Constraints

**minimum number of items**: the minimum number of items for this array is: `1`
