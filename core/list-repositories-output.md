# list-repositories output Schema

```txt
http://schema.nethserver.org/cluster/list-repositories-output.json
```

Output schema of the list-repositories action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                    |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-repositories-output.json](cluster/list-repositories-output.json "open original schema") |

## list-repositories output Type

`object[]` ([Details](list-repositories-output-items.md))

## list-repositories output Examples

```json
[
  {
    "name": "repository1",
    "url": "https://repository1.nethserver.org/",
    "testing": false,
    "status": true
  },
  {
    "name": "repository2",
    "url": "https://repository2.nethserver.org/",
    "testing": true,
    "status": true
  }
]
```
