# Untitled string in set-password-warning Schema

```txt
http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_template_content
```

The template content in base64 format, it can be a text file or an HTML file. It must respect python string.Template syntax. Required if mail\_template\_name is 'custom'

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [set-password-warning-input.json\*](cluster/set-password-warning-input.json "open original schema") |

## mail\_template\_content Type

`string`

## mail\_template\_content Constraints

**minimum length**: the minimum number of characters for this string is: `1`
