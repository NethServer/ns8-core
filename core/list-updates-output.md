# list-updates output Schema

```txt
http://schema.nethserver.org/cluster/list-updates-output.json
```

List updates output

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-updates-output.json](cluster/list-updates-output.json "open original schema") |

## list-updates output Type

`object[]` ([Details](list-updates-output-items.md))

## list-updates output Examples

```json
{
  "id": "dokuwiki1",
  "node": "1",
  "digest": "sha256:929465c177d9e40559b3fa838f2e429060ff7e8c3e4e1a12076b3304ad562982",
  "source": "ghcr.io/nethserver/dokuwiki",
  "version": "0.0.1-alpha1",
  "update": "0.0.2"
}
```
