# set-fqdn input Schema

```txt
http://schema.nethserver.org/node/set-fqdn-input.json
```

Input schema of the set-fqdn action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                             |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-fqdn-input.json](node/set-fqdn-input.json "open original schema") |

## set-fqdn input Type

`object` ([set-fqdn input](set-fqdn-input.md))

## set-fqdn input Examples

```json
{
  "hostname": "dn2",
  "domain": "example.org"
}
```

# set-fqdn input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                           |
| :-------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------- |
| [hostname](#hostname) | `string` | Required | cannot be null | [set-fqdn input](set-fqdn-input-properties-hostname.md "http://schema.nethserver.org/node/set-fqdn-input.json#/properties/hostname") |
| [domain](#domain)     | `string` | Required | cannot be null | [set-fqdn input](set-fqdn-input-properties-domain.md "http://schema.nethserver.org/node/set-fqdn-input.json#/properties/domain")     |

## hostname



`hostname`

* is required

* Type: `string`

* cannot be null

* defined in: [set-fqdn input](set-fqdn-input-properties-hostname.md "http://schema.nethserver.org/node/set-fqdn-input.json#/properties/hostname")

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



`domain`

* is required

* Type: `string`

* cannot be null

* defined in: [set-fqdn input](set-fqdn-input-properties-domain.md "http://schema.nethserver.org/node/set-fqdn-input.json#/properties/domain")

### domain Type

`string`

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")
