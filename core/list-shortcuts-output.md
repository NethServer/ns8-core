# list-shortcuts output Schema

```txt
http://schema.nethserver.org/cluster/list-shortcuts-output.json
```

Output schema of the list-shortcuts action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-shortcuts-output.json](cluster/list-shortcuts-output.json "open original schema") |

## list-shortcuts output Type

`object[]` ([Details](list-shortcuts-output-items.md))

## list-shortcuts output Examples

```json
[
  {
    "name": {
      "en": "Audit logs",
      "it": "Log di audit"
    },
    "description": {
      "en": "Show audit logs",
      "it": "Visualizza log di audit"
    },
    "tags": {
      "en": [
        "audit",
        "logs",
        "log",
        "operation"
      ],
      "it": [
        "audit",
        "logs",
        "log",
        "operazioni"
      ]
    },
    "path": "/logs",
    "source": "core"
  },
  {
    "name": {
      "en": "Change wiki URL",
      "it": "Cambia URL wiki"
    },
    "description": {
      "en": "Change the host name associated to the wiki",
      "it": "Cambia il nome host associato al wiki"
    },
    "tags": {
      "en": [
        "wiki",
        "dokuwiki",
        "host",
        "url"
      ],
      "it": [
        "wiki",
        "dokuwiki",
        "host",
        "url"
      ]
    },
    "path": "apps/dokuwiki1/edit",
    "source": "dokuwiki1"
  }
]
```
