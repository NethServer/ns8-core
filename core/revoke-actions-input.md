# revoke-actions input Schema

```txt
http://schema.nethserver.org/cluster/revoke-actions-input.json
```

Revoke permissions with matches. See also the grant-actions input schema.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [revoke-actions-input.json](cluster/revoke-actions-input.json "open original schema") |

## revoke-actions input Type

`object[]` ([Details](cluster-definitions-grant-assertion.md))

## revoke-actions input Examples

```json
[
  {
    "on": "cluster",
    "to": "reader",
    "action": "list-something"
  },
  {
    "action": "list-*",
    "to": "reader",
    "on": "cluster"
  }
]
```
