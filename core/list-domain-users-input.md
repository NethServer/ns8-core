# list-domain-users input Schema

```txt
http://schema.nethserver.org/cluster/list-domain-users-input.json
```

List users of a given accounts domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-domain-users-input.json](cluster/list-domain-users-input.json "open original schema") |

## list-domain-users input Type

`object` ([list-domain-users input](list-domain-users-input.md))

## list-domain-users input Examples

```json
{
  "domain": "dom.test"
}
```

# list-domain-users input Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                          |
| :---------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [domain](#domain) | `string` | Required | cannot be null | [list-domain-users input](list-domain-users-input-properties-domain-name.md "http://schema.nethserver.org/cluster/list-domain-users-input.json#/properties/domain") |

## domain



`domain`

* is required

* Type: `string` ([Domain name](list-domain-users-input-properties-domain-name.md))

* cannot be null

* defined in: [list-domain-users input](list-domain-users-input-properties-domain-name.md "http://schema.nethserver.org/cluster/list-domain-users-input.json#/properties/domain")

### domain Type

`string` ([Domain name](list-domain-users-input-properties-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`
