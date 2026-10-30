# add-external-domain input Schema

```txt
http://schema.nethserver.org/cluster/add-external-domain-input.json
```

Configure an external user domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-external-domain-input.json](cluster/add-external-domain-input.json "open original schema") |

## add-external-domain input Type

`object` ([add-external-domain input](add-external-domain-input.md))

any of

* not

  * [Protocol property is ldap](add-external-domain-input-anyof-0-protocol-property-is-ldap.md "check type definition")

* all of

  * [TCP service endpoint](add-external-domain-input-defs-tcp-service-endpoint.md "check type definition")

  * [LDAP domain properties](add-external-domain-input-defs-ldap-domain-properties.md "check type definition")

## add-external-domain input Examples

```json
{
  "domain": "example.com",
  "protocol": "ldap",
  "host": "18.19.20.21",
  "port": 636,
  "schema": "rfc2307",
  "bind_dn": "cn=ldapservice,dc=example,dc=com",
  "bind_password": "s3cret",
  "base_dn": "dc=example,dc=com",
  "tls": true,
  "tls_verify": true
}
```

# add-external-domain input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                                        |
| :-------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain)     | `string` | Required | cannot be null | [add-external-domain input](add-external-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/properties/domain")    |
| [protocol](#protocol) | `string` | Required | cannot be null | [add-external-domain input](add-external-domain-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/properties/protocol") |

## domain



`domain`

* is required

* Type: `string` ([User domain name](add-external-domain-input-properties-user-domain-name.md))

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/properties/domain")

### domain Type

`string` ([User domain name](add-external-domain-input-properties-user-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](add-external-domain-input-properties-provider-protocol.md))

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/properties/protocol")

### protocol Type

`string` ([Provider protocol](add-external-domain-input-properties-provider-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

# add-external-domain input Definitions

## Definitions group tcp-service-endpoint

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/tcp-service-endpoint"}
```

| Property      | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                            |
| :------------ | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [host](#host) | Merged    | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/tcp-service-endpoint/properties/host") |
| [port](#port) | `integer` | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/tcp-service-endpoint/properties/port") |

### host



`host`

* is required

* Type: `string` ([Details](add-external-domain-input-defs-tcp-service-endpoint-properties-host.md))

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/tcp-service-endpoint/properties/host")

#### host Type

`string` ([Details](add-external-domain-input-defs-tcp-service-endpoint-properties-host.md))

one (and only one) of

* [Untitled undefined type in add-external-domain input](add-external-domain-input-defs-tcp-service-endpoint-properties-host-oneof-0.md "check type definition")

* [Untitled undefined type in add-external-domain input](add-external-domain-input-defs-tcp-service-endpoint-properties-host-oneof-1.md "check type definition")

### port



`port`

* is required

* Type: `integer`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/tcp-service-endpoint/properties/port")

#### port Type

`integer`

#### port Constraints

**maximum**: the value of this number must smaller than or equal to: `65535`

**minimum**: the value of this number must greater than or equal to: `1`

## Definitions group additional-properties-of-ldap

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap"}
```

| Property                        | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                         |
| :------------------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [schema](#schema)               | `string`  | Optional | can be null    | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/schema") |
| [base_dn](#base_dn)             | `string`  | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/base_dn")             |
| [bind_dn](#bind_dn)             | `string`  | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_dn")             |
| [bind_password](#bind_password) | `string`  | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_password") |
| [tls](#tls)                     | `boolean` | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls")                     |
| [tls_verify](#tls_verify)       | `boolean` | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls_verify")       |

### schema

The LDAP schema is probed automatically if the value is `null` or the property is missing

`schema`

* is optional

* Type: `string` ([LDAP database schema](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md))

* can be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/schema")

#### schema Type

`string` ([LDAP database schema](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md))

#### schema Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value       | Explanation |
| :---------- | :---------- |
| `"ad"`      |             |
| `"rfc2307"` |             |

### base_dn

The LDAP base DN is probed automatically if the value is `""` (empty string)

`base_dn`

* is required

* Type: `string` ([Base DN](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md))

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/base_dn")

#### base_dn Type

`string` ([Base DN](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md))

### bind_dn



`bind_dn`

* is required

* Type: `string`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_dn")

#### bind_dn Type

`string`

#### bind_dn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### bind_password



`bind_password`

* is required

* Type: `string`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_password")

#### bind_password Type

`string`

#### bind_password Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### tls



`tls`

* is required

* Type: `boolean`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls")

#### tls Type

`boolean`

### tls_verify



`tls_verify`

* is required

* Type: `boolean`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls_verify")

#### tls_verify Type

`boolean`
