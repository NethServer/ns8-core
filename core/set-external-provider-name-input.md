# set-external-provider-name input Schema

```txt
http://schema.nethserver.org/cluster/set-external-provider-name-input.json
```

Set the UI display name of an external provider instance

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-external-provider-name-input.json](cluster/set-external-provider-name-input.json "open original schema") |

## set-external-provider-name input Type

`object` ([set-external-provider-name input](set-external-provider-name-input.md))

any of

* not

  * [Protocol property is ldap](set-external-provider-name-input-anyof-0-protocol-property-is-ldap.md "check type definition")

* [TCP service endpoint](set-external-provider-name-input-defs-tcp-service-endpoint.md "check type definition")

## set-external-provider-name input Examples

```json
{
  "domain": "example.com",
  "protocol": "ldap",
  "host": "ldap2.example.com",
  "port": 636,
  "ui_name": "Secondary"
}
```

# set-external-provider-name input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                   |
| :-------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ui\_name](#ui_name)  | `string` | Required | cannot be null | [set-external-provider-name input](set-external-provider-name-input-properties-provider-ui-display-name.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/properties/ui_name") |
| [domain](#domain)     | `string` | Required | cannot be null | [set-external-provider-name input](set-external-provider-name-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/properties/domain")          |
| [protocol](#protocol) | `string` | Required | cannot be null | [set-external-provider-name input](set-external-provider-name-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/properties/protocol")       |

## ui\_name



`ui_name`

* is required

* Type: `string` ([Provider UI display name](set-external-provider-name-input-properties-provider-ui-display-name.md))

* cannot be null

* defined in: [set-external-provider-name input](set-external-provider-name-input-properties-provider-ui-display-name.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/properties/ui_name")

### ui\_name Type

`string` ([Provider UI display name](set-external-provider-name-input-properties-provider-ui-display-name.md))

### ui\_name Constraints

**maximum length**: the maximum number of characters for this string is: `24`

## domain



`domain`

* is required

* Type: `string` ([User domain name](set-external-provider-name-input-properties-user-domain-name.md))

* cannot be null

* defined in: [set-external-provider-name input](set-external-provider-name-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/properties/domain")

### domain Type

`string` ([User domain name](set-external-provider-name-input-properties-user-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](set-external-provider-name-input-properties-provider-protocol.md))

* cannot be null

* defined in: [set-external-provider-name input](set-external-provider-name-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/properties/protocol")

### protocol Type

`string` ([Provider protocol](set-external-provider-name-input-properties-provider-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

# set-external-provider-name input Definitions

## Definitions group tcp-service-endpoint

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/$defs/tcp-service-endpoint"}
```

| Property      | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                 |
| :------------ | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [host](#host) | `string`  | Required | cannot be null | [set-external-provider-name input](set-external-provider-name-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/$defs/tcp-service-endpoint/properties/host") |
| [port](#port) | `integer` | Required | cannot be null | [set-external-provider-name input](set-external-provider-name-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/$defs/tcp-service-endpoint/properties/port") |

### host



`host`

* is required

* Type: `string`

* cannot be null

* defined in: [set-external-provider-name input](set-external-provider-name-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/$defs/tcp-service-endpoint/properties/host")

#### host Type

`string`

### port



`port`

* is required

* Type: `integer`

* cannot be null

* defined in: [set-external-provider-name input](set-external-provider-name-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/set-external-provider-name-input.json#/$defs/tcp-service-endpoint/properties/port")

#### port Type

`integer`

#### port Constraints

**maximum**: the value of this number must smaller than or equal to: `65535`

**minimum**: the value of this number must greater than or equal to: `1`
