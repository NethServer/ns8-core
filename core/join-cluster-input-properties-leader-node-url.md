# Leader node URL Schema

```txt
http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/url
```

Public URL of the leader API endpoint

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [join-cluster-input.json\*](cluster/join-cluster-input.json "open original schema") |

## url Type

`string` ([Leader node URL](join-cluster-input-properties-leader-node-url.md))

## url Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^https?://.+$
```

[try pattern](https://regexr.com/?expression=%5Ehttps%3F%3A%2F%2F.%2B%24 "try regular expression with regexr.com")
