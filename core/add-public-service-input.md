# add-public-service input Schema

```txt
http://schema.nethserver.org/node/add-public-service-input.json
```

Add firewall configuration for a new public service

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-public-service-input.json](node/add-public-service-input.json "open original schema") |

## add-public-service input Type

`object` ([add-public-service input](add-public-service-input.md))

## add-public-service input Examples

```json
{
  "service": "smtp",
  "ports": [
    "25/tcp"
  ]
}
```

# add-public-service input Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                          |
| :------------------ | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [service](#service) | `string` | Required | cannot be null | [add-public-service input](add-public-service-input-properties-service-name.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/service")               |
| [ports](#ports)     | `array`  | Required | cannot be null | [add-public-service input](add-public-service-input-properties-firewall-ports-configuration.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/ports") |

## service



`service`

*   is required

*   Type: `string` ([Service name](add-public-service-input-properties-service-name.md))

*   cannot be null

*   defined in: [add-public-service input](add-public-service-input-properties-service-name.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/service")

### service Type

`string` ([Service name](add-public-service-input-properties-service-name.md))

### service Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## ports



`ports`

*   is required

*   Type: `array` ([Firewall ports configuration](add-public-service-input-properties-firewall-ports-configuration.md))

*   cannot be null

*   defined in: [add-public-service input](add-public-service-input-properties-firewall-ports-configuration.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/ports")

### ports Type

`array` ([Firewall ports configuration](add-public-service-input-properties-firewall-ports-configuration.md))

### ports Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

### ports Examples

```json
"25/tcp"
```

```json
"21/udp"
```

```json
"9000-9010/tcp"
```
