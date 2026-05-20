# list-service-providers input Schema

```txt
http://schema.nethserver.org/agent/list-service-providers-input.json
```

Input schema of the basic list-service-providers action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-service-providers-input.json](agent/list-service-providers-input.json "open original schema") |

## list-service-providers input Type

`object` ([list-service-providers input](list-service-providers-input.md))

## list-service-providers input Examples

```json
{
  "service": "ldap",
  "transport": "tcp"
}
```

```json
{
  "service": "ldap",
  "transport": "tcp",
  "filter": {
    "node": "3"
  }
}
```

# list-service-providers input Properties

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                                                 |
| :---------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [service](#service)     | `string` | Required | cannot be null | [list-service-providers input](list-service-providers-input-properties-service-name.md "http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/service")         |
| [transport](#transport) | `string` | Optional | cannot be null | [list-service-providers input](list-service-providers-input-properties-transport-protocol.md "http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/transport") |
| [filter](#filter)       | `object` | Optional | cannot be null | [list-service-providers input](list-service-providers-input-properties-filter-clauses.md "http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/filter")        |

## service



`service`

* is required

* Type: `string` ([Service name](list-service-providers-input-properties-service-name.md))

* cannot be null

* defined in: [list-service-providers input](list-service-providers-input-properties-service-name.md "http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/service")

### service Type

`string` ([Service name](list-service-providers-input-properties-service-name.md))

### service Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## transport



`transport`

* is optional

* Type: `string` ([Transport protocol](list-service-providers-input-properties-transport-protocol.md))

* cannot be null

* defined in: [list-service-providers input](list-service-providers-input-properties-transport-protocol.md "http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/transport")

### transport Type

`string` ([Transport protocol](list-service-providers-input-properties-transport-protocol.md))

### transport Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### transport Examples

```json
"tcp"
```

```json
"udp"
```

```json
"http"
```

## filter

Return entries matching all the given clauses

`filter`

* is optional

* Type: `object` ([Filter clauses](list-service-providers-input-properties-filter-clauses.md))

* cannot be null

* defined in: [list-service-providers input](list-service-providers-input-properties-filter-clauses.md "http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/filter")

### filter Type

`object` ([Filter clauses](list-service-providers-input-properties-filter-clauses.md))
