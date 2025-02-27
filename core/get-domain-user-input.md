# get-domain-user input Schema

```txt
http://schema.nethserver.org/cluster/get-domain-user-input.json
```

Get the details of the given user

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-domain-user-input.json](cluster/get-domain-user-input.json "open original schema") |

## get-domain-user input Type

`object` ([get-domain-user input](get-domain-user-input.md))

## get-domain-user input Examples

```json
{
  "domain": "dom.test",
  "user": "Administrator"
}
```

# get-domain-user input Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                    |
| :---------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [domain](#domain) | `string` | Required | cannot be null | [get-domain-user input](get-domain-user-input-properties-domain-name.md "http://schema.nethserver.org/cluster/get-domain-user-input.json#/properties/domain") |
| [user](#user)     | `string` | Required | cannot be null | [get-domain-user input](get-domain-user-input-properties-user-name.md "http://schema.nethserver.org/cluster/get-domain-user-input.json#/properties/user")     |

## domain



`domain`

* is required

* Type: `string` ([Domain name](get-domain-user-input-properties-domain-name.md))

* cannot be null

* defined in: [get-domain-user input](get-domain-user-input-properties-domain-name.md "http://schema.nethserver.org/cluster/get-domain-user-input.json#/properties/domain")

### domain Type

`string` ([Domain name](get-domain-user-input-properties-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## user



`user`

* is required

* Type: `string` ([User name](get-domain-user-input-properties-user-name.md))

* cannot be null

* defined in: [get-domain-user input](get-domain-user-input-properties-user-name.md "http://schema.nethserver.org/cluster/get-domain-user-input.json#/properties/user")

### user Type

`string` ([User name](get-domain-user-input-properties-user-name.md))

### user Constraints

**minimum length**: the minimum number of characters for this string is: `1`
