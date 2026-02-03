# Untitled undefined type in list-backup-repositories output Schema

```txt
http://schema.nethserver.org/module/list-backup-repositories-output.json#/items/properties/is_generated_locally
```

Tells if the backup originates from the local cluster or from another cluster. The null value is returned if this information is missing completely, as it happens in old backups.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                   |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-backup-repositories-output.json\*](module/list-backup-repositories-output.json "open original schema") |

## is\_generated\_locally Type

`boolean`
