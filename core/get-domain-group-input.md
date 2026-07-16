# get-domain-group input Schema

```txt
http://schema.nethserver.org/cluster/get-domain-group-input.json
```

Get the details of the given group

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-domain-group-input.json](cluster/get-domain-group-input.json "open original schema") |

## get-domain-group input Type

`object` ([get-domain-group input](get-domain-group-input.md))

## get-domain-group input Examples

```json
{
  "domain": "dom.test",
  "group": "domain admins"
}
```

# get-domain-group input Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                       |
| :---------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain) | `string` | Required | cannot be null | [get-domain-group input](get-domain-group-input-properties-domain-name.md "http://schema.nethserver.org/cluster/get-domain-group-input.json#/properties/domain") |
| [group](#group)   | `string` | Required | cannot be null | [get-domain-group input](get-domain-group-input-properties-group-name.md "http://schema.nethserver.org/cluster/get-domain-group-input.json#/properties/group")   |

## domain



`domain`

* is required

* Type: `string` ([Domain name](get-domain-group-input-properties-domain-name.md))

* cannot be null

* defined in: [get-domain-group input](get-domain-group-input-properties-domain-name.md "http://schema.nethserver.org/cluster/get-domain-group-input.json#/properties/domain")

### domain Type

`string` ([Domain name](get-domain-group-input-properties-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## group



`group`

* is required

* Type: `string` ([Group name](get-domain-group-input-properties-group-name.md))

* cannot be null

* defined in: [get-domain-group input](get-domain-group-input-properties-group-name.md "http://schema.nethserver.org/cluster/get-domain-group-input.json#/properties/group")

### group Type

`string` ([Group name](get-domain-group-input-properties-group-name.md))

### group Constraints

**minimum length**: the minimum number of characters for this string is: `1`
