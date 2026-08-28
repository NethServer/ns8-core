# Backup object Schema

```txt
http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-backups-output.json\*](cluster/list-backups-output.json "open original schema") |

## backup-item Type

`object` ([Backup object](list-backups-output-defs-backup-object.md))

# backup-item Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                              |
| :------------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)                        | `integer` | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-backup-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/id")                                                  |
| [instances](#instances)          | `array`   | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-module-instances.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/instances")                                    |
| [schedule](#schedule)            | `string`  | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule")                             |
| [schedule\_hint](#schedule_hint) | `object`  | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule_hint")                 |
| [retention](#retention)          | `integer` | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-backup-retention.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/retention")                                    |
| [name](#name)                    | `string`  | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-backup-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/name")                                              |
| [enabled](#enabled)              | `boolean` | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/enabled") |

## id



`id`

* is required

* Type: `integer` ([Backup ID](list-backups-output-defs-backup-object-properties-backup-id.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-backup-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/id")

### id Type

`integer` ([Backup ID](list-backups-output-defs-backup-object-properties-backup-id.md))

### id Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## instances

Identifiers of module instances included in the backup

`instances`

* is required

* Type: `object[]` ([Backup of a module instance](list-backups-output-defs-backup-of-a-module-instance.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-module-instances.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/instances")

### instances Type

`object[]` ([Backup of a module instance](list-backups-output-defs-backup-of-a-module-instance.md))

## schedule

The string format must be compatible with Calendar Events, described in `systemd.timer` man page

`schedule`

* is required

* Type: `string` ([Time schedule expression](list-backups-output-defs-backup-object-properties-time-schedule-expression.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule")

### schedule Type

`string` ([Time schedule expression](list-backups-output-defs-backup-object-properties-time-schedule-expression.md))

### schedule Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `1`

## schedule\_hint

Store arbitrary object to ease parsing of `schedule` value

`schedule_hint`

* is required

* Type: `object` ([Schedule expression hint for UI](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule_hint")

### schedule\_hint Type

`object` ([Schedule expression hint for UI](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md))

## retention

Number of backup snapshots to store before automatic deletion

`retention`

* is required

* Type: `integer` ([Backup retention](list-backups-output-defs-backup-object-properties-backup-retention.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-backup-retention.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/retention")

### retention Type

`integer` ([Backup retention](list-backups-output-defs-backup-object-properties-backup-retention.md))

### retention Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## name



`name`

* is required

* Type: `string` ([Backup name](list-backups-output-defs-backup-object-properties-backup-name.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-backup-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/name")

### name Type

`string` ([Backup name](list-backups-output-defs-backup-object-properties-backup-name.md))

### name Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## enabled



`enabled`

* is required

* Type: `boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/enabled")

### enabled Type

`boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))
