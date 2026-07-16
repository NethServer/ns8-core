# list-service-providers output Schema

```txt
http://schema.nethserver.org/agent/list-service-providers-output.json
```

Output schema of the basic list-service-providers action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-service-providers-output.json](agent/list-service-providers-output.json "open original schema") |

## list-service-providers output Type

`array` ([list-service-providers output](list-service-providers-output.md))

## list-service-providers output Examples

```json
[
  {
    "port": "143",
    "host": "10.5.4.1",
    "node": "1",
    "user_domain": "dp.nethserver.net",
    "module_uuid": "8d257122-0a7f-49c7-a620-08961a68cfa0",
    "module_id": "mail1",
    "ui_name": null
  }
]
```
