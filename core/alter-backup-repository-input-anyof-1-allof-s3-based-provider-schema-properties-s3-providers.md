# S3 providers Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/anyOf/1/allOf/0/properties/provider
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## provider Type

unknown ([S3 providers](alter-backup-repository-input-anyof-1-allof-s3-based-provider-schema-properties-s3-providers.md))

## provider Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value          | Explanation |
| :------------- | :---------- |
| `"aws"`        |             |
| `"generic-s3"` |             |
