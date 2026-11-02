# validate-leader-fqdn input Schema

```txt
http://schema.nethserver.org/node/validate-leader-fqdn-input.json
```

Input schema of the validate-leader-fqdn action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                     |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [validate-leader-fqdn-input.json](node/validate-leader-fqdn-input.json "open original schema") |

## validate-leader-fqdn input Type

`object` ([validate-leader-fqdn input](validate-leader-fqdn-input.md))

## validate-leader-fqdn input Examples

```json
{
  "fqdn": "dn2.example.org"
}
```

# validate-leader-fqdn input Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                       |
| :------------ | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [fqdn](#fqdn) | `string` | Required | cannot be null | [validate-leader-fqdn input](validate-leader-fqdn-input-properties-fqdn.md "http://schema.nethserver.org/node/validate-leader-fqdn-input.json#/properties/fqdn") |

## fqdn



`fqdn`

* is required

* Type: `string`

* cannot be null

* defined in: [validate-leader-fqdn input](validate-leader-fqdn-input-properties-fqdn.md "http://schema.nethserver.org/node/validate-leader-fqdn-input.json#/properties/fqdn")

### fqdn Type

`string`

### fqdn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")
