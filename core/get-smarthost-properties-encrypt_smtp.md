# encrypt\_smtp Schema

```txt
http://schema.nethserver.org/cluster/get-smarthost.json#/properties/encrypt_smtp
```

Enable or disable the tls encryption with the smtp server

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [get-smarthost.json\*](cluster/get-smarthost.json "open original schema") |

## encrypt\_smtp Type

`string` ([encrypt\_smtp](get-smarthost-properties-encrypt_smtp.md))

## encrypt\_smtp Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value        | Explanation |
| :----------- | :---------- |
| `"none"`     |             |
| `"starttls"` |             |
| `"tls"`      |             |
