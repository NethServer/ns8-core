# Untitled string in set-password-warning Schema

```txt
http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_subject
```

The email subject. It must respect python string.Template syntax. Required if mail\_template\_name is 'custom'

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [set-password-warning-input.json\*](cluster/set-password-warning-input.json "open original schema") |

## mail\_subject Type

`string`

## mail\_subject Constraints

**minimum length**: the minimum number of characters for this string is: `1`
