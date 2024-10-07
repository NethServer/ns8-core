# list-backups output Schema

```txt
http://schema.nethserver.org/cluster/list-backups-output.json
```

Get a list of backup configurations

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-backups-output.json](cluster/list-backups-output.json "open original schema") |

## list-backups output Type

`object` ([list-backups output](list-backups-output.md))

## list-backups output Examples

```json
{
  "backups": [
    {
      "id": 1,
      "name": "Daily",
      "repository": "48ce000a-79b7-5fe6-8558-177fd70c27b4",
      "schedule": "daily",
      "schedule_hint": {},
      "retention": "7d",
      "instances": [
        {
          "module_id": "dokuwiki1",
          "ui_name": "",
          "repository_path": "dokuwiki1@2f72561e-89b2-4cdc-b4e4-425ca23bbec9",
          "status": {
            "total_size": 4053660,
            "total_file_count": 21744,
            "start": 1640097808,
            "end": 1640097815,
            "success": true
          }
        }
      ],
      "enabled": true
    }
  ],
  "unconfigured_instances": [
    "nextcloud3"
  ]
}
```

# list-backups output Properties

| Property                                           | Type    | Required | Nullable       | Defined by                                                                                                                                                                                |
| :------------------------------------------------- | :------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [backups](#backups)                                | `array` | Required | cannot be null | [list-backups output](list-backups-output-properties-backups.md "http://schema.nethserver.org/cluster/list-backups-output.json#/properties/backups")                                      |
| [unconfigured\_instances](#unconfigured_instances) | `array` | Required | cannot be null | [list-backups output](list-backups-output-properties-unconfigured-module-instances.md "http://schema.nethserver.org/cluster/list-backups-output.json#/properties/unconfigured_instances") |

## backups

List of configured backup objects

`backups`

* is required

* Type: `object[]` ([Backup object](list-backups-output-defs-backup-object.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-properties-backups.md "http://schema.nethserver.org/cluster/list-backups-output.json#/properties/backups")

### backups Type

`object[]` ([Backup object](list-backups-output-defs-backup-object.md))

## unconfigured\_instances

List of module instances that are not included in any backup

`unconfigured_instances`

* is required

* Type: `object[]` ([Unconfigured module instance](list-backups-output-defs-unconfigured-module-instance.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-properties-unconfigured-module-instances.md "http://schema.nethserver.org/cluster/list-backups-output.json#/properties/unconfigured_instances")

### unconfigured\_instances Type

`object[]` ([Unconfigured module instance](list-backups-output-defs-unconfigured-module-instance.md))

# list-backups output Definitions

## Definitions group uinstance-item

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item"}
```

| Property             | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                         |
| :------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)            | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/id")           |
| [ui\_name](#ui_name) | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/ui_name") |

### id



`id`

* is optional

* Type: `string` ([Module ID](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/id")

#### id Type

`string` ([Module ID](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md))

#### id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### ui\_name



`ui_name`

* is optional

* Type: `string` ([Module UI name](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/ui_name")

#### ui\_name Type

`string` ([Module UI name](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md))

## Definitions group instance-item

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item"}
```

| Property                             | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                |
| :----------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [module\_id](#module_id)             | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-module_id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/module_id")             |
| [ui\_name](#ui_name-1)               | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-ui_name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/ui_name")                 |
| [repository\_path](#repository_path) | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-repository_path.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/repository_path") |
| [status](#status)                    | `object` | Optional | can be null    | [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/status")   |

### module\_id



`module_id`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-module_id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/module_id")

#### module\_id Type

`string`

#### module\_id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### ui\_name



`ui_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-ui_name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/ui_name")

#### ui\_name Type

`string`

### repository\_path



`repository_path`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-repository_path.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/repository_path")

#### repository\_path Type

`string`

#### repository\_path Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### status



`status`

* is optional

* Type: `object` ([Last backup run status](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md))

* can be null

* defined in: [list-backups output](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/instance-item/properties/status")

#### status Type

`object` ([Last backup run status](list-backups-output-defs-backup-of-a-module-instance-properties-last-backup-run-status.md))

## Definitions group backup-status

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-status"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group backup-item

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item"}
```

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                              |
| :------------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id-1)                      | `integer` | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-backup-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/id")                                                  |
| [instances](#instances)          | `array`   | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-module-instances.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/instances")                                    |
| [schedule](#schedule)            | `string`  | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule")                             |
| [schedule\_hint](#schedule_hint) | `object`  | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule_hint")                 |
| [retention](#retention)          | `integer` | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-backup-retention.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/retention")                                    |
| [name](#name)                    | `string`  | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-backup-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/name")                                              |
| [enabled](#enabled)              | `boolean` | Required | cannot be null | [list-backups output](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/enabled") |

### id



`id`

* is required

* Type: `integer` ([Backup ID](list-backups-output-defs-backup-object-properties-backup-id.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-backup-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/id")

#### id Type

`integer` ([Backup ID](list-backups-output-defs-backup-object-properties-backup-id.md))

#### id Constraints

**minimum**: the value of this number must greater than or equal to: `1`

### instances

Identifiers of module instances included in the backup

`instances`

* is required

* Type: `object[]` ([Backup of a module instance](list-backups-output-defs-backup-of-a-module-instance.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-module-instances.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/instances")

#### instances Type

`object[]` ([Backup of a module instance](list-backups-output-defs-backup-of-a-module-instance.md))

### schedule

The string format must be compatible with Calendar Events, described in `systemd.timer` man page

`schedule`

* is required

* Type: `string` ([Time schedule expression](list-backups-output-defs-backup-object-properties-time-schedule-expression.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-time-schedule-expression.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule")

#### schedule Type

`string` ([Time schedule expression](list-backups-output-defs-backup-object-properties-time-schedule-expression.md))

#### schedule Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `1`

### schedule\_hint

Store arbitrary object to ease parsing of `schedule` value

`schedule_hint`

* is required

* Type: `object` ([Schedule expression hint for UI](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/schedule_hint")

#### schedule\_hint Type

`object` ([Schedule expression hint for UI](list-backups-output-defs-backup-object-properties-schedule-expression-hint-for-ui.md))

### retention

Number of backup snapshots to store before automatic deletion

`retention`

* is required

* Type: `integer` ([Backup retention](list-backups-output-defs-backup-object-properties-backup-retention.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-backup-retention.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/retention")

#### retention Type

`integer` ([Backup retention](list-backups-output-defs-backup-object-properties-backup-retention.md))

#### retention Constraints

**minimum**: the value of this number must greater than or equal to: `1`

### name



`name`

* is required

* Type: `string` ([Backup name](list-backups-output-defs-backup-object-properties-backup-name.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-backup-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/name")

#### name Type

`string` ([Backup name](list-backups-output-defs-backup-object-properties-backup-name.md))

#### name Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### enabled



`enabled`

* is required

* Type: `boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/backup-item/properties/enabled")

#### enabled Type

`boolean` ([Enabled/disabled status. If \`false\` the backup is stopped.](list-backups-output-defs-backup-object-properties-enableddisabled-status-if-false-the-backup-is-stopped.md))
