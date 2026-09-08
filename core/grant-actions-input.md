# grant-actions input Schema

```txt
http://schema.nethserver.org/cluster/grant-actions-input.json
```

Set permissions with a list of grant assertions

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [grant-actions-input.json](cluster/grant-actions-input.json "open original schema") |

## grant-actions input Type

`object[]` ([Details](cluster-definitions-grant-assertion.md))

## grant-actions input Examples

```json
[
  {
    "action": "*",
    "to": "owner",
    "on": "cluster"
  },
  {
    "action": "list-*",
    "to": "*",
    "on": "cluster"
  }
]
```
