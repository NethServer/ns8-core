# Untitled string in add-custom-zone input Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1/properties/destinations/items
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-custom-zone-input.json\*](node/add-custom-zone-input.json "open original schema") |

## items Type

`string`

## items Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[0-9]{1,3}(\.[0-9]{1,3}){3}(/[0-9]{1,2})?(:[0-9]{1,5}(-[0-9]{1,5})?/(tcp|udp))?$
```

[try pattern](https://regexr.com/?expression=%5E%5B0-9%5D%7B1%2C3%7D\(%5C.%5B0-9%5D%7B1%2C3%7D\)%7B3%7D\(%2F%5B0-9%5D%7B1%2C2%7D\)%3F\(%3A%5B0-9%5D%7B1%2C5%7D\(-%5B0-9%5D%7B1%2C5%7D\)%3F%2F\(tcp%7Cudp\)\)%3F%24 "try regular expression with regexr.com")
