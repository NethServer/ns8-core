# run-backup input Schema

```txt
http://schema.nethserver.org/agent/run-backup-input.json
```

Run the given backup immediately

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [run-backup-input.json](agent/run-backup-input.json "open original schema") |

## run-backup input Type

`object` ([run-backup input](run-backup-input-1.md))

## run-backup input Examples

```json
{
  "id": 5
}
```

# run-backup input Properties

| Property  | Type      | Required | Nullable       | Defined by                                                                                                                               |
| :-------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id) | `integer` | Required | cannot be null | [run-backup input](run-backup-input-1-properties-backup-id.md "http://schema.nethserver.org/agent/run-backup-input.json#/properties/id") |

## id



`id`

* is required

* Type: `integer` ([Backup ID](run-backup-input-1-properties-backup-id.md))

* cannot be null

* defined in: [run-backup input](run-backup-input-1-properties-backup-id.md "http://schema.nethserver.org/agent/run-backup-input.json#/properties/id")

### id Type

`integer` ([Backup ID](run-backup-input-1-properties-backup-id.md))

### id Constraints

**minimum**: the value of this number must greater than or equal to: `1`
