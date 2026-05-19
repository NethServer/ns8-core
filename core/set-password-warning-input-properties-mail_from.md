# Untitled string in set-password-warning Schema

```txt
http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_from
```

The sender email address

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [set-password-warning-input.json\*](cluster/set-password-warning-input.json "open original schema") |

## mail\_from Type

`string`

## mail\_from Constraints

**email**: the string must be an email address, according to [RFC 5322, section 3.4.1](https://tools.ietf.org/html/rfc5322 "check the specification")
