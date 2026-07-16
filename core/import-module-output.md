# import-module output Schema

```txt
http://schema.nethserver.org/cluster/import-module-output.json
```

Output schema of the import-module action. It has the same format of the add-module action output

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [import-module-output.json](cluster/import-module-output.json "open original schema") |

## import-module output Type

`object` ([import-module output](import-module-output.md))

## import-module output Examples

```json
{
  "credentials": [
    "dokuwiki1",
    "16541c262d4c8-b1e8-4e95-a31b-8c3ed239203f"
  ],
  "port": 20031,
  "address": "10.5.4.1",
  "task": "module/dokuwiki1/b75964f0-766c-4869-bacd-941171753487"
}
```
