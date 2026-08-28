# Time schedule expression Schema

```txt
http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule
```

The string format must be compatible with Calendar Events, described in `systemd.timer` man page

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-backups-output.json\*](cluster/list-backups-output.json "open original schema") |

## schedule Type

`string` ([Time schedule expression](list-backups-output-defs-backup-object-properties-time-schedule-expression.md))

## schedule Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `1`
