# remove-backup-repository input Schema

```txt
http://schema.nethserver.org/cluster/remove-backup-repository-input.json
```

Remove a backup repository and any related backup object

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-backup-repository-input.json](cluster/remove-backup-repository-input.json "open original schema") |

## remove-backup-repository input Type

`object` ([remove-backup-repository input](remove-backup-repository-input.md))

## remove-backup-repository input Examples

```json
{
  "id": "183f7ae2-7649-5bb5-8742-9ee24d058b8b"
}
```

# remove-backup-repository input Properties

| Property  | Type     | Required | Nullable       | Defined by                                                                                                                                                                             |
| :-------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id) | `string` | Required | cannot be null | [remove-backup-repository input](remove-backup-repository-input-properties-repository-id.md "http://schema.nethserver.org/cluster/remove-backup-repository-input.json#/properties/id") |

## id



`id`

* is required

* Type: `string` ([Repository ID](remove-backup-repository-input-properties-repository-id.md))

* cannot be null

* defined in: [remove-backup-repository input](remove-backup-repository-input-properties-repository-id.md "http://schema.nethserver.org/cluster/remove-backup-repository-input.json#/properties/id")

### id Type

`string` ([Repository ID](remove-backup-repository-input-properties-repository-id.md))

### id Constraints

**minimum length**: the minimum number of characters for this string is: `1`
