# Grant assertions list Schema

```txt
http://schema.nethserver.org/cluster/alter-user-input.json#/properties/grant
```

A list of initial roles on the matching objects

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [alter-user-input.json\*](cluster/alter-user-input.json "open original schema") |

## grant Type

`object[]` ([Grant object](cluster-definitions-grant-object.md))

## grant Constraints

**unique items**: all items in this array must be unique. Duplicates are not allowed.
