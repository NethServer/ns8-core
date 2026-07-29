# list-users output Schema

```txt
http://schema.nethserver.org/cluster/list-users-output.json
```

Output schema of the list-users action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-users-output.json](cluster/list-users-output.json "open original schema") |

## list-users output Type

`object[]` ([Details](list-users-output-items.md))

## list-users output Examples

```json
[
  {
    "display_name": "Admin",
    "user": "admin"
  },
  {
    "display_name": "People",
    "user": "people"
  }
]
```
