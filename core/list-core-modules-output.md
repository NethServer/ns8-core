# list-core-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-core-modules-output.json
```

List core modules output

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                    |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-core-modules-output.json](cluster/list-core-modules-output.json "open original schema") |

## list-core-modules output Type

`object[]` ([Details](list-core-modules-output-items.md))

## list-core-modules output Examples

```json
[
  {
    "name": "core",
    "instances": [
      {
        "id": "core",
        "version": "updates_from_repo",
        "update": "",
        "node_id": "1",
        "node_ui_name": "node1"
      }
    ]
  },
  {
    "name": "promtail",
    "instances": [
      {
        "id": "promtail1",
        "version": "latest",
        "update": "",
        "node_id": "1",
        "node_ui_name": "node1"
      }
    ]
  },
  {
    "name": "traefik",
    "instances": [
      {
        "id": "traefik1",
        "version": "0.0.1",
        "update": "",
        "node_id": "1",
        "node_ui_name": "node1"
      }
    ]
  },
  {
    "name": "loki",
    "instances": [
      {
        "id": "loki1",
        "version": "latest",
        "update": "",
        "node_id": "1",
        "node_ui_name": "node1"
      },
      {
        "id": "loki2",
        "version": "0.0.1-alpha1",
        "update": "0.0.1",
        "node_id": "2",
        "node_ui_name": "node2"
      }
    ]
  },
  {
    "name": "ldapproxy",
    "instances": [
      {
        "id": "ldapproxy1",
        "version": "latest",
        "update": "",
        "node_id": "1",
        "node_ui_name": "node1"
      }
    ]
  }
]
```
