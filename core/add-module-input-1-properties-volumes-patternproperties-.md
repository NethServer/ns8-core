# Untitled string in add-module input Schema

```txt
http://schema.nethserver.org/cluster/add-module-input.json#/properties/volumes/patternProperties/[^/]+
```

Absolute filesystem directory path under /

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-module-input.json\*](cluster/add-module-input.json "open original schema") |

## ]+ Type

`string`

## ]+ Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")
