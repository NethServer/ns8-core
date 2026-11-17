# Backup retention Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/retention
```

Number of backup snapshots to store before automatic deletion

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [alter-backup-input.json\*](cluster/alter-backup-input.json "open original schema") |

## retention Type

`integer` ([Backup retention](alter-backup-input-properties-backup-retention.md))

## retention Constraints

**minimum**: the value of this number must greater than or equal to: `1`
