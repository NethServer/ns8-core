# read-backup-repositories output Schema

```txt
http://schema.nethserver.org/cluster/read-backup-repositories-output.json
```

Read the content of all backup repositories

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                  |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [read-backup-repositories-output.json](cluster/read-backup-repositories-output.json "open original schema") |

## read-backup-repositories output Type

unknown ([read-backup-repositories output](read-backup-repositories-output.md))

## read-backup-repositories output Examples

```json
[
  {
    "name": "dokuwiki",
    "path": "dokuwiki/dokuwiki1@cc7335d7-4d67-408c-8c35-42257667e51b",
    "uuid": "cc7335d7-4d67-408c-8c35-42257667e51b",
    "timestamp": 1644403745,
    "repository_id": "e181c936-1bc3-5032-a809-d8b0551eebe9",
    "repository_name": "B2 repo",
    "repository_provider": "backblaze",
    "repository_url": "b2:giacomons8",
    "installed_instance": "dokuwiki4",
    "installed_instance_ui_name": "My Dokuwiki"
  }
]
```
