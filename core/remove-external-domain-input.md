# remove-external-domain input Schema

```txt
http://schema.nethserver.org/cluster/remove-external-domain-input.json
```

Remove an external user domain and all its providers

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-external-domain-input.json](cluster/remove-external-domain-input.json "open original schema") |

## remove-external-domain input Type

`object` ([remove-external-domain input](remove-external-domain-input.md))

## remove-external-domain input Examples

```json
{
  "domain": "example.com"
}
```

# remove-external-domain input Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                                              |
| :---------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain) | `string` | Required | cannot be null | [remove-external-domain input](remove-external-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/remove-external-domain-input.json#/properties/domain") |

## domain



`domain`

* is required

* Type: `string` ([User domain name](remove-external-domain-input-properties-user-domain-name.md))

* cannot be null

* defined in: [remove-external-domain input](remove-external-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/remove-external-domain-input.json#/properties/domain")

### domain Type

`string` ([User domain name](remove-external-domain-input-properties-user-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`
