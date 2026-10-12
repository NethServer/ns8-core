# Time schedule expression Schema

```txt
http://schema.nethserver.org/cluster/add-backup-input.json#/properties/schedule
```

The string format must be compatible with Calendar Events, described in `systemd.timer` man page

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-backup-input.json\*](cluster/add-backup-input.json "open original schema") |

## schedule Type

`string` ([Time schedule expression](add-backup-input-properties-time-schedule-expression.md))

## schedule Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `1`
