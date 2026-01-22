# Untitled string in list-alerts output Schema

```txt
http://schema.nethserver.org/cluster/list-alerts.json#/$defs/alertmanager-v2-object/properties/startsAt
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-alerts.json\*](cluster/list-alerts.json "open original schema") |

## startsAt Type

`string`

## startsAt Constraints

**date time**: the string must be a date time string, according to [RFC 3339, section 5.6](https://tools.ietf.org/html/rfc3339 "check the specification")
