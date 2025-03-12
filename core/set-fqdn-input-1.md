# set-fqdn input Schema

```txt
http://schema.nethserver.org/cluster/set-fqdn-input.json
```

Input schema of the set-fqdn action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-fqdn-input.json](cluster/set-fqdn-input.json "open original schema") |

## set-fqdn input Type

`object` ([set-fqdn input](set-fqdn-input-1.md))

## set-fqdn input Examples

```json
{
  "node": 3,
  "hostname": "dn2",
  "domain": "example.org",
  "reachability_check": true
}
```

# set-fqdn input Properties

| Property                                   | Type      | Required | Nullable       | Defined by                                                                                                                                                    |
| :----------------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [node](#node)                              | `integer` | Required | cannot be null | [set-fqdn input](set-fqdn-input-1-properties-node.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/node")                             |
| [reachability\_check](#reachability_check) | `boolean` | Optional | cannot be null | [set-fqdn input](set-fqdn-input-1-properties-reachability_check.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/reachability_check") |
| [hostname](#hostname)                      | `string`  | Required | cannot be null | [set-fqdn input](set-fqdn-input-1-properties-hostname.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/hostname")                     |
| [domain](#domain)                          | `string`  | Required | cannot be null | [set-fqdn input](set-fqdn-input-1-properties-domain.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/domain")                         |

## node

The identifier of the cluster node.

`node`

* is required

* Type: `integer`

* cannot be null

* defined in: [set-fqdn input](set-fqdn-input-1-properties-node.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/node")

### node Type

`integer`

## reachability\_check

Check if the new FQDN can be correctly resolved and is reachable by HTTP from other cluster nodes.

`reachability_check`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [set-fqdn input](set-fqdn-input-1-properties-reachability_check.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/reachability_check")

### reachability\_check Type

`boolean`

## hostname

The hostname part of the FQDN for the node.

`hostname`

* is required

* Type: `string`

* cannot be null

* defined in: [set-fqdn input](set-fqdn-input-1-properties-hostname.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/hostname")

### hostname Type

`string`

### hostname Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[^\.]+$
```

[try pattern](https://regexr.com/?expression=%5E%5B%5E%5C.%5D%2B%24 "try regular expression with regexr.com")

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")

## domain

The domain part of the FQDN for the node.

`domain`

* is required

* Type: `string`

* cannot be null

* defined in: [set-fqdn input](set-fqdn-input-1-properties-domain.md "http://schema.nethserver.org/cluster/set-fqdn-input.json#/properties/domain")

### domain Type

`string`

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")
