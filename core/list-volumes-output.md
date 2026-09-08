# List module volumes Schema

```txt
http://schema.nethserver.org/agent/list-volumes-output.json
```

A list of Podman volume names of the current module

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-volumes-output.json](agent/list-volumes-output.json "open original schema") |

## List module volumes Type

`string[]`

## List module volumes Examples

```json
[
  {
    "name": "dokuwiki-data",
    "created": "2021-07-19 10:19:45.528366456 +0200 CEST"
  }
]
```
