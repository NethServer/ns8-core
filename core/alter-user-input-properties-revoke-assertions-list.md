# Revoke assertions list Schema

```txt
http://schema.nethserver.org/cluster/alter-user-input.json#/properties/revoke
```

A list of revoke roles on the matching objects

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [alter-user-input.json\*](cluster/alter-user-input.json "open original schema") |

## revoke Type

`object[]` ([Grant object](cluster-definitions-grant-object.md))

## revoke Constraints

**unique items**: all items in this array must be unique. Duplicates are not allowed.
