# Untitled string in validate-leader-fqdn input Schema

```txt
http://schema.nethserver.org/node/validate-leader-fqdn-input.json#/properties/fqdn
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                       |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [validate-leader-fqdn-input.json\*](node/validate-leader-fqdn-input.json "open original schema") |

## fqdn Type

`string`

## fqdn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")
