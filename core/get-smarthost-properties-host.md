# host Schema

```txt
http://schema.nethserver.org/cluster/get-smarthost.json#/properties/host
```

Host name for the smarthost, like 'wiki.nethserver.org'

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [get-smarthost.json\*](cluster/get-smarthost.json "open original schema") |

## host Type

`string` ([host](get-smarthost-properties-host.md))

## host Constraints

**(international) hostname**: the string must be an (IDN) hostname, according to [RFC 5890, section 2.3.2.3](https://tools.ietf.org/html/rfc5890 "check the specification")
