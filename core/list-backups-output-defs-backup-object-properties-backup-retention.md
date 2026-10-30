# Backup retention Schema

```txt
http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/retention
```

Number of backup snapshots to store before automatic deletion

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-backups-output.json\*](cluster/list-backups-output.json "open original schema") |

## retention Type

`integer` ([Backup retention](list-backups-output-defs-backup-object-properties-backup-retention.md))

## retention Constraints

**minimum**: the value of this number must greater than or equal to: `1`
