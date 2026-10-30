# add-backup output Schema

```txt
http://schema.nethserver.org/cluster/add-backup-output.json
```

The add-backup action returns the backup ID of the added item

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-backup-output.json](cluster/add-backup-output.json "open original schema") |

## add-backup output Type

`integer` ([add-backup output](add-backup-output.md))

## add-backup output Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## add-backup output Examples

```json
1
```
