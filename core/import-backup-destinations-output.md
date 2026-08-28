# import-backup-destinations output Schema

```txt
http://schema.nethserver.org/cluster/import-backup-destinations-output.json
```

Import backup destination configurations from raw cluster backup data

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [import-backup-destinations-output.json](cluster/import-backup-destinations-output.json "open original schema") |

## import-backup-destinations output Type

`object` ([import-backup-destinations output](import-backup-destinations-output.md))

## import-backup-destinations output Examples

```json
{
  "imported_destinations": 2,
  "skipped_destinations": 1
}
```

# import-backup-destinations output Properties

| Property                                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                 |
| :----------------------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [imported\_destinations](#imported_destinations) | `integer` | Required | cannot be null | [import-backup-destinations output](import-backup-destinations-output-properties-imported_destinations.md "http://schema.nethserver.org/cluster/import-backup-destinations-output.json#/properties/imported_destinations") |
| [skipped\_destinations](#skipped_destinations)   | `integer` | Required | cannot be null | [import-backup-destinations output](import-backup-destinations-output-properties-skipped_destinations.md "http://schema.nethserver.org/cluster/import-backup-destinations-output.json#/properties/skipped_destinations")   |

## imported\_destinations

Number of backup destinations successfully imported

`imported_destinations`

* is required

* Type: `integer`

* cannot be null

* defined in: [import-backup-destinations output](import-backup-destinations-output-properties-imported_destinations.md "http://schema.nethserver.org/cluster/import-backup-destinations-output.json#/properties/imported_destinations")

### imported\_destinations Type

`integer`

### imported\_destinations Constraints

**minimum**: the value of this number must greater than or equal to: `0`

## skipped\_destinations

Number of backup destinations not imported because already configured in the cluster DB

`skipped_destinations`

* is required

* Type: `integer`

* cannot be null

* defined in: [import-backup-destinations output](import-backup-destinations-output-properties-skipped_destinations.md "http://schema.nethserver.org/cluster/import-backup-destinations-output.json#/properties/skipped_destinations")

### skipped\_destinations Type

`integer`

### skipped\_destinations Constraints

**minimum**: the value of this number must greater than or equal to: `0`
