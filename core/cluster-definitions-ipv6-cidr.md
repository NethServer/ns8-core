# IPv6 CIDR Schema

```txt
http://schema.nethserver.org/cluster.json#/definitions/ipv6-cidr
```

IPv6 address or network in CIDR notation

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [cluster.json\*](cluster.json "open original schema") |

## ipv6-cidr Type

`string` ([IPv6 CIDR](cluster-definitions-ipv6-cidr.md))

## ipv6-cidr Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[0-9a-fA-F:]+(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?$
```

[try pattern](https://regexr.com/?expression=%5E%5B0-9a-fA-F%3A%5D%2B\(%5C%2F\(%5B0-9%5D%7C%5B1-9%5D%5B0-9%5D%7C1%5B0-1%5D%5B0-9%5D%7C12%5B0-8%5D\)\)%3F%24 "try regular expression with regexr.com")

## ipv6-cidr Examples

```json
"::1"
```

```json
"fd12:3456:789a::/48"
```
