# remove-backup input Schema

```txt
http://schema.nethserver.org/cluster/remove-backup-input.json
```

Remove a backup object

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-backup-input.json](cluster/remove-backup-input.json "open original schema") |

## remove-backup input Type

`object` ([remove-backup input](remove-backup-input.md))

## remove-backup input Examples

```json
{
  "id": 7
}
```

# remove-backup input Properties

| Property  | Type      | Required | Nullable       | Defined by                                                                                                                                        |
| :-------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| [id](#id) | `integer` | Required | cannot be null | [remove-backup input](remove-backup-input-properties-backup-id.md "http://schema.nethserver.org/cluster/remove-backup-input.json#/properties/id") |

## id



`id`

* is required

* Type: `integer` ([Backup ID](remove-backup-input-properties-backup-id.md))

* cannot be null

* defined in: [remove-backup input](remove-backup-input-properties-backup-id.md "http://schema.nethserver.org/cluster/remove-backup-input.json#/properties/id")

### id Type

`integer` ([Backup ID](remove-backup-input-properties-backup-id.md))
