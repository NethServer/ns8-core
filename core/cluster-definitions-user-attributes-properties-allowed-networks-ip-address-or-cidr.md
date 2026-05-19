# IP address or CIDR Schema

```txt
http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/allowed_networks/items
```

IPv4 or IPv6 address or network in CIDR notation

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [cluster.json\*](cluster.json "open original schema") |

## items Type

merged type ([IP address or CIDR](cluster-definitions-user-attributes-properties-allowed-networks-ip-address-or-cidr.md))

any of

* [IPv4 CIDR](cluster-definitions-ip-address-or-cidr-anyof-ipv4-cidr.md "check type definition")

* [IPv6 CIDR](cluster-definitions-ip-address-or-cidr-anyof-ipv6-cidr.md "check type definition")
