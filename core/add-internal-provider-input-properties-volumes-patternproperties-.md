# Untitled string in add-internal-provider input Schema

```txt
http://schema.nethserver.org/cluster/add-internal-provider-input.json#/properties/volumes/patternProperties/[^/]+
```

Absolute filesystem directory path under /

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-internal-provider-input.json\*](cluster/add-internal-provider-input.json "open original schema") |

## ]+ Type

`string`

## ]+ Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^/
```

[try pattern](https://regexr.com/?expression=%5E%2F "try regular expression with regexr.com")
