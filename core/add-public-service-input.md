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
    "587/tcp"
  ]
}
```

```json
{
  "service": "smtp",
  "replace_ports": true,
  "ports": [
    "465/tcp",
    "587/tcp",
    "25/tcp"
  ]
}
```

# add-public-service input Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                          |
| :------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [service](#service)              | `string`  | Required | cannot be null | [add-public-service input](add-public-service-input-properties-service-name.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/service")               |
| [replace\_ports](#replace_ports) | `boolean` | Optional | cannot be null | [add-public-service input](add-public-service-input-properties-replace_ports.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/replace_ports")        |
| [ports](#ports)                  | `array`   | Required | cannot be null | [add-public-service input](add-public-service-input-properties-firewall-ports-configuration.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/ports") |

## service



`service`

* is required

* Type: `string` ([Service name](add-public-service-input-properties-service-name.md))

* cannot be null

* defined in: [add-public-service input](add-public-service-input-properties-service-name.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/service")

### service Type

`string` ([Service name](add-public-service-input-properties-service-name.md))

### service Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## replace\_ports

If true, the provided list of ports replaces any existing configuration for the service. If false or omitted, the ports are added to the existing set.

`replace_ports`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [add-public-service input](add-public-service-input-properties-replace_ports.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/replace_ports")

### replace\_ports Type

`boolean`

## ports



`ports`

* is required

* Type: `array` ([Firewall ports configuration](add-public-service-input-properties-firewall-ports-configuration.md))

* cannot be null

* defined in: [add-public-service input](add-public-service-input-properties-firewall-ports-configuration.md "http://schema.nethserver.org/node/add-public-service-input.json#/properties/ports")

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
