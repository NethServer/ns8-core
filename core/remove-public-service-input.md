# remove-public-service input Schema

```txt
http://schema.nethserver.org/node/remove-public-service-input.json
```

Remove firewall configuration for the given public service

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                       |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-public-service-input.json](node/remove-public-service-input.json "open original schema") |

## remove-public-service input Type

`object` ([remove-public-service input](remove-public-service-input.md))

## remove-public-service input Examples

```json
{
  "service": "smtp"
}
```

# remove-public-service input Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                     |
| :------------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [service](#service) | `string` | Required | cannot be null | [remove-public-service input](remove-public-service-input-properties-service-name.md "http://schema.nethserver.org/node/remove-public-service-input.json#/properties/service") |

## service



`service`

* is required

* Type: `string` ([Service name](remove-public-service-input-properties-service-name.md))

* cannot be null

* defined in: [remove-public-service input](remove-public-service-input-properties-service-name.md "http://schema.nethserver.org/node/remove-public-service-input.json#/properties/service")

### service Type

`string` ([Service name](remove-public-service-input-properties-service-name.md))

### service Constraints

**minimum length**: the minimum number of characters for this string is: `1`
