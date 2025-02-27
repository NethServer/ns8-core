# remove-custom-zone input Schema

```txt
http://schema.nethserver.org/node/remove-custom-zone-input.json
```

Remove firewall configuration for the given zone

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-custom-zone-input.json](node/remove-custom-zone-input.json "open original schema") |

## remove-custom-zone input Type

`object` ([remove-custom-zone input](remove-custom-zone-input.md))

## remove-custom-zone input Examples

```json
{
  "zone": "myzone"
}
```

# remove-custom-zone input Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                         |
| :------------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [service](#service) | `string` | Optional | cannot be null | [remove-custom-zone input](remove-custom-zone-input-properties-zone-name.md "http://schema.nethserver.org/node/remove-custom-zone-input.json#/properties/service") |

## service



`service`

* is optional

* Type: `string` ([Zone name](remove-custom-zone-input-properties-zone-name.md))

* cannot be null

* defined in: [remove-custom-zone input](remove-custom-zone-input-properties-zone-name.md "http://schema.nethserver.org/node/remove-custom-zone-input.json#/properties/service")

### service Type

`string` ([Zone name](remove-custom-zone-input-properties-zone-name.md))

### service Constraints

**minimum length**: the minimum number of characters for this string is: `1`
