# Untitled string in update-modules input Schema

```txt
http://schema.nethserver.org/cluster/update-modules-input.json#/properties/instances/items
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [update-modules-input.json\*](cluster/update-modules-input.json "open original schema") |

## items Type

`string`

## items Constraints

**minimum length**: the minimum number of characters for this string is: `2`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^.+[0-9]+$
```

[try pattern](https://regexr.com/?expression=%5E.%2B%5B0-9%5D%2B%24 "try regular expression with regexr.com")
