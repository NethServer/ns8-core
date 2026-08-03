# alter-backup input Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-input.json
```

Configure a new backup instance

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-backup-input.json](cluster/alter-backup-input.json "open original schema") |

## alter-backup input Type

`object` ([alter-backup input](alter-backup-input.md))

## alter-backup input Examples

```json
{
  "id": 5,
  "name": "BB daily",
  "schedule": "daily",
  "schedule_hint": {},
  "retention": "7d",
  "instances": [
    "dokuwiki1",
    "nextcloud1"
  ],
  "enabled": true
}
```

# alter-backup input Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                      |
| :------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)                        | `integer` | Required | cannot be null | [alter-backup input](alter-backup-input-properties-backup-id.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/id")                                                  |
| [instances](#instances)          | `array`   | Required | cannot be null | [alter-backup input](alter-backup-input-properties-module-instances.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/instances")                                    |
| [schedule](#schedule)            | `string`  | Required | cannot be null | [alter-backup input](alter-backup-input-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/schedule")                             |
| [schedule\_hint](#schedule_hint) | `object`  | Optional | cannot be null | [alter-backup input](alter-backup-input-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/schedule_hint")                 |
| [retention](#retention)          | `integer` | Required | cannot be null | [alter-backup input](alter-backup-input-properties-backup-retention.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/retention")                                    |
| [name](#name)                    | `string`  | Required | cannot be null | [alter-backup input](alter-backup-input-properties-backup-name.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/name")                                              |
| [enabled](#enabled)              | `boolean` | Required | cannot be null | [alter-backup input](alter-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/enabled") |

## id



`id`

* is required

* Type: `integer` ([Backup ID](alter-backup-input-properties-backup-id.md))

* cannot be null

* defined in: [alter-backup input](alter-backup-input-properties-backup-id.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/id")

### id Type

`integer` ([Backup ID](alter-backup-input-properties-backup-id.md))

### id Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## instances

Identifiers of module instances included in the backup

`instances`

* is required

* Type: `string[]`

* cannot be null

* defined in: [alter-backup input](alter-backup-input-properties-module-instances.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/instances")

### instances Type

`string[]`

## schedule

The string format must be compatible with Calendar Events, described in `systemd.timer` man page

`schedule`

* is required

* Type: `string` ([Time schedule expression](alter-backup-input-properties-time-schedule-expression.md))

* cannot be null

* defined in: [alter-backup input](alter-backup-input-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/schedule")

### schedule Type

`string` ([Time schedule expression](alter-backup-input-properties-time-schedule-expression.md))

### schedule Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `1`

## schedule\_hint

Store arbitrary object to ease parsing of `schedule` value

`schedule_hint`

* is optional

* Type: `object` ([Schedule expression hint for UI](alter-backup-input-properties-schedule-expression-hint-for-ui.md))

* cannot be null

* defined in: [alter-backup input](alter-backup-input-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/schedule_hint")

### schedule\_hint Type

`object` ([Schedule expression hint for UI](alter-backup-input-properties-schedule-expression-hint-for-ui.md))

## retention

Number of backup snapshots to store before automatic deletion

`retention`

* is required

* Type: `integer` ([Backup retention](alter-backup-input-properties-backup-retention.md))

* cannot be null

* defined in: [alter-backup input](alter-backup-input-properties-backup-retention.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/retention")

### retention Type

`integer` ([Backup retention](alter-backup-input-properties-backup-retention.md))

### retention Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## name



`name`

* is required

* Type: `string` ([Backup name](alter-backup-input-properties-backup-name.md))

* cannot be null

* defined in: [alter-backup input](alter-backup-input-properties-backup-name.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/name")

### name Type

`string` ([Backup name](alter-backup-input-properties-backup-name.md))

### name Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## enabled



`enabled`

* is required

* Type: `boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](alter-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))

* cannot be null

* defined in: [alter-backup input](alter-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/alter-backup-input.json#/properties/enabled")

### enabled Type

`boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](alter-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))
