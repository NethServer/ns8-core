# IPv4 CIDR Schema

```txt
http://schema.nethserver.org/cluster/join-node-input.json#/properties/network
```

IPv4 with netmask in CIDR notation

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [join-node-input.json\*](cluster/join-node-input.json "open original schema") |

## network Type

`string` ([IPv4 CIDR](cluster-definitions-ipv4-cidr.md))

## network Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$
```

[try pattern](https://regexr.com/?expression=%5E\(%5B0-9%5D%7B1%2C3%7D%5C.\)%7B3%7D%5B0-9%5D%7B1%2C3%7D\(%5C%2F\(%5B0-9%5D%7C%5B1-2%5D%5B0-9%5D%7C3%5B0-2%5D\)\)%3F%24 "try regular expression with regexr.com")

## network Examples

```json
"10.5.4.0/24"
```

```json
"192.168.73.0/24"
```
