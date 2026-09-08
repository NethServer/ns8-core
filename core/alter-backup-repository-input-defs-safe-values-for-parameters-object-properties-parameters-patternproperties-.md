# Untitled undefined type in alter-backup-repository input Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/safe_parameters/properties/parameters/patternProperties/.*
```

Strings without forbidden control chars

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## .\* Type

unknown

## .\* Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[^\n\r\t]*$
```

[try pattern](https://regexr.com/?expression=%5E%5B%5E%5Cn%5Cr%5Ct%5D*%24 "try regular expression with regexr.com")
