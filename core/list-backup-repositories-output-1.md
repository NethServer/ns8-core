# list-backup-repositories output Schema

```txt
http://schema.nethserver.org/cluster/list-backup-repositories-output.json
```

Get the list of available backup destinations and the status of cluster backup password

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                  |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-backup-repositories-output.json](cluster/list-backup-repositories-output.json "open original schema") |

## list-backup-repositories output Type

`object` ([list-backup-repositories output](list-backup-repositories-output-1.md))

## list-backup-repositories output Examples

```json
{
  "repositories": [
    {
      "id": "48ce000a-79b7-5fe6-8558-177fd70c27b4",
      "provider": "backblaze",
      "name": "BackBlaze repo1",
      "url": "b2:backupex1",
      "password": "d59a90ec7ad2b2967257a7a308c82c96ac006efd138254bc1e58c8ea07c18400",
      "parameters": {
        "b2_account_id": "xxxxxxxxxxxxxx",
        "b2_account_key": "yyyyyyyyyyyyyyyyyyyyyy"
      }
    }
  ],
  "password_exists": true
}
```
