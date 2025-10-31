# configure-backup input Schema

```txt
http://schema.nethserver.org/agent/configure-backup-input.json
```

Input schema of the basic configure-backup action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [configure-backup-input.json](agent/configure-backup-input.json "open original schema") |

## configure-backup input Type

`object` ([configure-backup input](configure-backup-input.md))

any of

* [Add or alter the backup configuration](configure-backup-input-anyof-add-or-alter-the-backup-configuration.md "check type definition")

* [Remove the backup configuration](configure-backup-input-anyof-remove-the-backup-configuration.md "check type definition")

## configure-backup input Examples

```json
{
  "optype": "add",
  "id": 5,
  "name": "BB daily",
  "schedule": "daily",
  "enabled": true
}
```

# configure-backup input Properties

| Property              | Type      | Required | Nullable       | Defined by                                                                                                                                                        |
| :-------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [optype](#optype)     | `string`  | Optional | cannot be null | [configure-backup input](configure-backup-input-properties-operation-type.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/optype") |
| [id](#id)             | `integer` | Optional | cannot be null | [configure-backup input](configure-backup-input-properties-backup-id.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/id")          |
| [name](#name)         | `string`  | Optional | cannot be null | [configure-backup input](configure-backup-input-properties-backup-name.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/name")      |
| [schedule](#schedule) | `string`  | Optional | cannot be null | [configure-backup input](configure-backup-input-properties-schedule.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/schedule")     |
| [enabled](#enabled)   | `boolean` | Optional | cannot be null | [configure-backup input](configure-backup-input-properties-enabled.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/enabled")       |

## optype



`optype`

* is optional

* Type: `string` ([Operation type](configure-backup-input-properties-operation-type.md))

* cannot be null

* defined in: [configure-backup input](configure-backup-input-properties-operation-type.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/optype")

### optype Type

`string` ([Operation type](configure-backup-input-properties-operation-type.md))

## id



`id`

* is optional

* Type: `integer` ([Backup ID](configure-backup-input-properties-backup-id.md))

* cannot be null

* defined in: [configure-backup input](configure-backup-input-properties-backup-id.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/id")

### id Type

`integer` ([Backup ID](configure-backup-input-properties-backup-id.md))

### id Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## name



`name`

* is optional

* Type: `string` ([Backup name](configure-backup-input-properties-backup-name.md))

* cannot be null

* defined in: [configure-backup input](configure-backup-input-properties-backup-name.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/name")

### name Type

`string` ([Backup name](configure-backup-input-properties-backup-name.md))

## schedule



`schedule`

* is optional

* Type: `string`

* cannot be null

* defined in: [configure-backup input](configure-backup-input-properties-schedule.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/schedule")

### schedule Type

`string`

## enabled



`enabled`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [configure-backup input](configure-backup-input-properties-enabled.md "http://schema.nethserver.org/agent/configure-backup-input.json#/properties/enabled")

### enabled Type

`boolean`
