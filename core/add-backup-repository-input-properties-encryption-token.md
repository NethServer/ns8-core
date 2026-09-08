# Encryption token Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-input.json#/properties/password
```

Select the password for restic encryption. If this is empty the API will generate a random password

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-backup-repository-input.json\*](cluster/add-backup-repository-input.json "open original schema") |

## password Type

`string` ([Encryption token](add-backup-repository-input-properties-encryption-token.md))
