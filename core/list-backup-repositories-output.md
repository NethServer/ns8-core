# list-backup-repositories output Schema

```txt
http://schema.nethserver.org/module/list-backup-repositories-output.json
```

Return a list of the module's Restic backup repositories. The format is the same of cluster/read-backup-repositories.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                 |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-backup-repositories-output.json](module/list-backup-repositories-output.json "open original schema") |

## list-backup-repositories output Type

`object[]` ([Details](list-backup-repositories-output-items.md))

## list-backup-repositories output Examples

```json
[
  {
    "module_id": "loki1",
    "module_ui_name": "My Loki",
    "node_fqdn": "rl1.dp.nethserver.net",
    "path": "loki/35f45b73-f81e-467b-b622-96ec3b7fec19",
    "name": "loki",
    "uuid": "35f45b73-f81e-467b-b622-96ec3b7fec19",
    "timestamp": 1721405723,
    "repository_id": "14030a59-a4e6-54cc-b8ea-cd5f97fe44c8",
    "repository_name": "BackBlaze repo1",
    "repository_provider": "backblaze",
    "repository_url": "b2:ns8-davidep",
    "installed_instance": "loki1",
    "installed_instance_ui_name": "My Loki",
    "is_generated_locally": true
  }
]
```
