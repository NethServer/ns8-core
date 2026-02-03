# Leader IP address Schema

```txt
http://schema.nethserver.org/cluster/join-node-input.json#/properties/leader_ip_address
```

VPN IP address of the leader node

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [join-node-input.json\*](cluster/join-node-input.json "open original schema") |

## leader\_ip\_address Type

`string` ([Leader IP address](join-node-input-properties-leader-ip-address.md))

## leader\_ip\_address Constraints

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")
