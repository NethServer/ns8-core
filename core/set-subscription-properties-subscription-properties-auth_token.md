# Untitled string in set-subscription Schema

```txt
http://schema.nethserver.org/cluster/set-subscription.json#/properties/subscription/properties/auth_token
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [set-subscription.json\*](cluster/set-subscription.json "open original schema") |

## auth\_token Type

`string`

## auth\_token Constraints

**maximum length**: the maximum number of characters for this string is: `128`

**minimum length**: the minimum number of characters for this string is: `32`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^\S+$
```

[try pattern](https://regexr.com/?expression=%5E%5CS%2B%24 "try regular expression with regexr.com")
