# Untitled string in set-fqdn input Schema

```txt
http://schema.nethserver.org/node/set-fqdn-input.json#/properties/domain
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                               |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [set-fqdn-input.json\*](node/set-fqdn-input.json "open original schema") |

## domain Type

`string`

## domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")
