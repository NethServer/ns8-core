# remove-internal-domain input Schema

```txt
http://schema.nethserver.org/cluster/remove-internal-domain-input.json
```

Remove an internal user domain and all its providers

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-internal-domain-input.json](cluster/remove-internal-domain-input.json "open original schema") |

## remove-internal-domain input Type

`object` ([remove-internal-domain input](remove-internal-domain-input.md))

## remove-internal-domain input Examples

```json
{
  "domain": "example.com"
}
```

# remove-internal-domain input Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                                              |
| :---------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain) | `string` | Required | cannot be null | [remove-internal-domain input](remove-internal-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/remove-internal-domain-input.json#/properties/domain") |

## domain



`domain`

* is required

* Type: `string` ([User domain name](remove-internal-domain-input-properties-user-domain-name.md))

* cannot be null

* defined in: [remove-internal-domain input](remove-internal-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/remove-internal-domain-input.json#/properties/domain")

### domain Type

`string` ([User domain name](remove-internal-domain-input-properties-user-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`
