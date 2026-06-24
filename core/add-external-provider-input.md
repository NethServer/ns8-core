# add-external-provider input Schema

```txt
http://schema.nethserver.org/cluster/add-external-provider-input.json
```

Add a provider to an already configured external user domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-external-provider-input.json](cluster/add-external-provider-input.json "open original schema") |

## add-external-provider input Type

`object` ([add-external-provider input](add-external-provider-input.md))

any of

* not

  * [Protocol property is ldap](add-external-provider-input-anyof-0-protocol-property-is-ldap.md "check type definition")

* [TCP service endpoint](add-external-provider-input-defs-tcp-service-endpoint.md "check type definition")

## add-external-provider input Examples

```json
{
  "domain": "example.com",
  "protocol": "ldap",
  "host": "ldap2.example.com",
  "port": 636
}
```

# add-external-provider input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                                              |
| :-------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain)     | `string` | Required | cannot be null | [add-external-provider input](add-external-provider-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/properties/domain")    |
| [protocol](#protocol) | `string` | Required | cannot be null | [add-external-provider input](add-external-provider-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/properties/protocol") |

## domain



`domain`

* is required

* Type: `string` ([User domain name](add-external-provider-input-properties-user-domain-name.md))

* cannot be null

* defined in: [add-external-provider input](add-external-provider-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/properties/domain")

### domain Type

`string` ([User domain name](add-external-provider-input-properties-user-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](add-external-provider-input-properties-provider-protocol.md))

* cannot be null

* defined in: [add-external-provider input](add-external-provider-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/properties/protocol")

### protocol Type

`string` ([Provider protocol](add-external-provider-input-properties-provider-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

# add-external-provider input Definitions

## Definitions group tcp-service-endpoint

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint"}
```

| Property      | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                  |
| :------------ | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [host](#host) | `string`  | Required | cannot be null | [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/host") |
| [port](#port) | `integer` | Required | cannot be null | [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/port") |

### host



`host`

* is required

* Type: `string`

* cannot be null

* defined in: [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/host")

#### host Type

`string`

### port



`port`

* is required

* Type: `integer`

* cannot be null

* defined in: [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/port")

#### port Type

`integer`

#### port Constraints

**maximum**: the value of this number must smaller than or equal to: `65535`

**minimum**: the value of this number must greater than or equal to: `1`
