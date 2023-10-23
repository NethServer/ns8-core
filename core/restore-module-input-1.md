# restore-module input Schema

```txt
http://schema.nethserver.org/cluster/restore-module-input.json
```

Input schema of the restore-module action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [restore-module-input.json](cluster/restore-module-input.json "open original schema") |

## restore-module input Type

`object[]` ([Details](restore-module-input-1-items.md))

## restore-module input Examples

```json
[
  {
    "repository": "48ce000a-79b7-5fe6-8558-177fd70c27b4",
    "path": "dokuwiki/dokuwiki1@f5d24fcd-819c-4b1d-98ad-a1b2ebcee8cf",
    "snapshot": "",
    "node": 1
  }
]
```
