# Untitled undefined type in alter-backup-repository input Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/rclone_parameters/properties/rclone_conf_secret
```

Basic INI-like string validation

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## rclone\_conf\_secret Type

`string`

## rclone\_conf\_secret Constraints

**maximum length**: the maximum number of characters for this string is: `100000`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^(\[[^\]]+\]\n)?(\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*[^\n]*\n)*$
```

[try pattern](https://regexr.com/?expression=%5E\(%5C%5B%5B%5E%5C%5D%5D%2B%5C%5D%5Cn\)%3F\(%5Cs*\(%5BA-Za-z_%5D%5BA-Za-z0-9_%5D*\)%5Cs*%3D%5Cs*%5B%5E%5Cn%5D*%5Cn\)*%24 "try regular expression with regexr.com")
