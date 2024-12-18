# Endpoint address Schema

```txt
http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_address
```

Host name or IP address where the new leader node can be reached by other nodes

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [promote-node-input.json\*](cluster/promote-node-input.json "open original schema") |

## endpoint\_address Type

`string` ([Endpoint address](promote-node-input-properties-endpoint-address.md))

## endpoint\_address Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")
