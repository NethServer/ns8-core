# list-favorites output Schema

```txt
http://schema.nethserver.org/cluster/list-favorites-output.json
```

Output schema of the list-favorites action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-favorites-output.json](cluster/list-favorites-output.json "open original schema") |

## list-favorites output Type

`object[]` ([Details](list-favorites-output-items.md))

## list-favorites output Examples

```json
[
  {
    "id": "dokuwiki1",
    "module": "dokuwiki"
  },
  {
    "id": "nextcloud1",
    "module": "nextcloud"
  }
]
```
