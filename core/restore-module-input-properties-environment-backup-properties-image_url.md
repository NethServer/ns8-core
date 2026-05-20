# Untitled string in restore-module input Schema

```txt
http://schema.nethserver.org/module/restore-module-input.json#/properties/environment/properties/IMAGE_URL
```

URL of the module root image

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [restore-module-input.json\*](module/restore-module-input.json "open original schema") |

## IMAGE\_URL Type

`string`

## IMAGE\_URL Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## IMAGE\_URL Examples

```json
"ghcr.io/nethserver/mymodule:v2.3.2"
```
