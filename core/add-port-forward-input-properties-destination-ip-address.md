# Destination IP address Schema

```txt
http://schema.nethserver.org/node/add-port-forward-input.json#/properties/destination_ip
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                               |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-port-forward-input.json\*](node/add-port-forward-input.json "open original schema") |

## destination\_ip Type

`string` ([Destination IP address](add-port-forward-input-properties-destination-ip-address.md))

## destination\_ip Constraints

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")
