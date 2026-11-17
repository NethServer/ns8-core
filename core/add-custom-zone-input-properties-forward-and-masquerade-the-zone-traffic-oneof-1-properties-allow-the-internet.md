# Allow the Internet Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1/properties/internet
```

Private networks stay closed unless a destination lists them

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-custom-zone-input.json\*](node/add-custom-zone-input.json "open original schema") |

## internet Type

`boolean` ([Allow the Internet](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1-properties-allow-the-internet.md))

## internet Default Value

The default value is:

```json
false
```
