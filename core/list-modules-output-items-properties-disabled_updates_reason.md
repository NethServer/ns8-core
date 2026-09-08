# Untitled undefined type in list-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/disabled_updates_reason
```

A non-empty strings indicates that updates are disabled for some reason. The string value identifies the reason.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-modules-output.json\*](cluster/list-modules-output.json "open original schema") |

## disabled\_updates\_reason Type

unknown

## disabled\_updates\_reason Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value             | Explanation |
| :---------------- | :---------- |
| `""`              |             |
| `"ns7_migration"` |             |
