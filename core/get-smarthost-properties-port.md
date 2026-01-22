# port Schema

```txt
http://schema.nethserver.org/cluster/get-smarthost.json#/properties/port
```

Port number of the smtp smarthost

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [get-smarthost.json\*](cluster/get-smarthost.json "open original schema") |

## port Type

`integer` ([port](get-smarthost-properties-port.md))

## port Constraints

**maximum**: the value of this number must smaller than or equal to: `65535`

**minimum**: the value of this number must greater than or equal to: `1`
