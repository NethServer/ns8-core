# Allowed networks Schema

```txt
http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/allowed_networks
```

List of IPv4/IPv6 addresses or CIDR networks the user is allowed to connect from. An empty list removes any restriction.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [cluster.json\*](cluster.json "open original schema") |

## allowed\_networks Type

an array of merged types ([IP address or CIDR](cluster-definitions-user-attributes-properties-allowed-networks-ip-address-or-cidr.md))

## allowed\_networks Constraints

**unique items**: all items in this array must be unique. Duplicates are not allowed.
