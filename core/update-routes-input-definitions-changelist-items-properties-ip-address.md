# IP address Schema

```txt
http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/ip_address
```

IP address to add or remove. It should be local to the node.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [update-routes-input.json\*](cluster/update-routes-input.json "open original schema") |

## ip\_address Type

`string` ([IP address](update-routes-input-definitions-changelist-items-properties-ip-address.md))

## ip\_address Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")
