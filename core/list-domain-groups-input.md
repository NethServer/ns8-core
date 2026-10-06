# list-domain-groups input Schema

```txt
http://schema.nethserver.org/cluster/list-domain-groups-input.json
```

List groups of a given accounts domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-domain-groups-input.json](cluster/list-domain-groups-input.json "open original schema") |

## list-domain-groups input Type

`object` ([list-domain-groups input](list-domain-groups-input.md))

## list-domain-groups input Examples

```json
{
  "domain": "dom.test"
}
```

# list-domain-groups input Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                             |
| :---------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain) | `string` | Required | cannot be null | [list-domain-groups input](list-domain-groups-input-properties-domain-name.md "http://schema.nethserver.org/cluster/list-domain-groups-input.json#/properties/domain") |

## domain



`domain`

* is required

* Type: `string` ([Domain name](list-domain-groups-input-properties-domain-name.md))

* cannot be null

* defined in: [list-domain-groups input](list-domain-groups-input-properties-domain-name.md "http://schema.nethserver.org/cluster/list-domain-groups-input.json#/properties/domain")

### domain Type

`string` ([Domain name](list-domain-groups-input-properties-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`
