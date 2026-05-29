# Untitled string in List loki instances Schema

```txt
http://schema.nethserver.org/cluster/list-loki-instances-output.json#/properties/instances/items/properties/active_from
```

The ISO 8601 date-time when the Loki instance was activated.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-loki-instances-output.json\*](cluster/list-loki-instances-output.json "open original schema") |

## active\_from Type

`string`

## active\_from Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")
