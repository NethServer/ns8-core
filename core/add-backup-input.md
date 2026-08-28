# add-backup input Schema

```txt
http://schema.nethserver.org/cluster/add-backup-input.json
```

Configure a new backup instance

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-backup-input.json](cluster/add-backup-input.json "open original schema") |

## add-backup input Type

`object` ([add-backup input](add-backup-input.md))

## add-backup input Examples

```json
{
  "name": "BB daily",
  "repository": "6a2a6208-7d4f-5915-9419-e4f1faaa1c76",
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

# add-backup input Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                |
| :------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [instances](#instances)          | `array`   | Required | cannot be null | [add-backup input](add-backup-input-properties-module-instances.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/instances")                                    |
| [repository](#repository)        | `string`  | Required | cannot be null | [add-backup input](add-backup-input-properties-repository-uuid-reference.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/repository")                          |
| [schedule](#schedule)            | `string`  | Required | cannot be null | [add-backup input](add-backup-input-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/schedule")                             |
| [schedule\_hint](#schedule_hint) | `object`  | Optional | cannot be null | [add-backup input](add-backup-input-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/schedule_hint")                 |
| [retention](#retention)          | `integer` | Required | cannot be null | [add-backup input](add-backup-input-properties-backup-retention.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/retention")                                    |
| [name](#name)                    | `string`  | Required | cannot be null | [add-backup input](add-backup-input-properties-backup-name.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/name")                                              |
| [enabled](#enabled)              | `boolean` | Required | cannot be null | [add-backup input](add-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/enabled") |

## instances

Identifiers of module instances included in the backup

`instances`

* is required

* Type: `string[]`

* cannot be null

* defined in: [add-backup input](add-backup-input-properties-module-instances.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/instances")

### instances Type

`string[]`

## repository



`repository`

* is required

* Type: `string` ([Repository UUID reference](add-backup-input-properties-repository-uuid-reference.md))

* cannot be null

* defined in: [add-backup input](add-backup-input-properties-repository-uuid-reference.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/repository")

### repository Type

`string` ([Repository UUID reference](add-backup-input-properties-repository-uuid-reference.md))

### repository Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
```

[try pattern](https://regexr.com/?expression=%5E%5B0-9a-f%5D%7B8%7D-%5B0-9a-f%5D%7B4%7D-%5B0-9a-f%5D%7B4%7D-%5B0-9a-f%5D%7B4%7D-%5B0-9a-f%5D%7B12%7D%24 "try regular expression with regexr.com")

## schedule

The string format must be compatible with Calendar Events, described in `systemd.timer` man page

`schedule`

* is required

* Type: `string` ([Time schedule expression](add-backup-input-properties-time-schedule-expression.md))

* cannot be null

* defined in: [add-backup input](add-backup-input-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/schedule")

### schedule Type

`string` ([Time schedule expression](add-backup-input-properties-time-schedule-expression.md))

### schedule Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `1`

## schedule\_hint

Store arbitrary object to ease parsing of `schedule` value

`schedule_hint`

* is optional

* Type: `object` ([Schedule expression hint for UI](add-backup-input-properties-schedule-expression-hint-for-ui.md))

* cannot be null

* defined in: [add-backup input](add-backup-input-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/schedule_hint")

### schedule\_hint Type

`object` ([Schedule expression hint for UI](add-backup-input-properties-schedule-expression-hint-for-ui.md))

## retention

Number of backup snapshots to store before automatic deletion

`retention`

* is required

* Type: `integer` ([Backup retention](add-backup-input-properties-backup-retention.md))

* cannot be null

* defined in: [add-backup input](add-backup-input-properties-backup-retention.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/retention")

### retention Type

`integer` ([Backup retention](add-backup-input-properties-backup-retention.md))

### retention Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## name

If set to empty, a default name is set (e.g. `Backup 3`)

`name`

* is required

* Type: `string` ([Backup name](add-backup-input-properties-backup-name.md))

* cannot be null

* defined in: [add-backup input](add-backup-input-properties-backup-name.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/name")

### name Type

`string` ([Backup name](add-backup-input-properties-backup-name.md))

## enabled



`enabled`

* is required

* Type: `boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](add-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))

* cannot be null

* defined in: [add-backup input](add-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/add-backup-input.json#/properties/enabled")

### enabled Type

`boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](add-backup-input-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))
