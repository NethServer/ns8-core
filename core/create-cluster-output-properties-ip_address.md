# Untitled string in create-cluster output Schema

```txt
http://schema.nethserver.org/cluster/create-cluster-output.json#/properties/ip_address
```

The VPN IP address of the leader node

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [create-cluster-output.json\*](cluster/create-cluster-output.json "open original schema") |

## ip\_address Type

`string`

## ip\_address Constraints

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")
