# get-password-warning Schema

```txt
http://schema.nethserver.org/cluster/get-password-warning-input.json
```

Input schema of the get-password-warning action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-password-warning-input.json](cluster/get-password-warning-input.json "open original schema") |

## get-password-warning Type

`object` ([get-password-warning](get-password-warning-input.md))

## get-password-warning Examples

```json
{
  "domain": "mydom.test"
}
```

# get-password-warning Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                        |
| :---------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain) | `string` | Required | cannot be null | [get-password-warning](get-password-warning-input-properties-domain.md "http://schema.nethserver.org/cluster/get-password-warning-input.json#/properties/domain") |

## domain

A user domain name

`domain`

* is required

* Type: `string`

* cannot be null

* defined in: [get-password-warning](get-password-warning-input-properties-domain.md "http://schema.nethserver.org/cluster/get-password-warning-input.json#/properties/domain")

### domain Type

`string`

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`
