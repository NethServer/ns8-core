# alter-external-domain input Schema

```txt
http://schema.nethserver.org/cluster/alter-external-domain-input.json
```

Configure an external user domain

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-external-domain-input.json](cluster/alter-external-domain-input.json "open original schema") |

## alter-external-domain input Type

`object` ([alter-external-domain input](alter-external-domain-input.md))

## alter-external-domain input Examples

```json
{
  "domain": "example.com",
  "protocol": "ldap",
  "bind_dn": "cn=ldapservice,dc=example,dc=com",
  "bind_password": "s3cret",
  "tls": true,
  "tls_verify": true
}
```

# alter-external-domain input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                                              |
| :-------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain)     | `string` | Required | cannot be null | [alter-external-domain input](alter-external-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/properties/domain")    |
| [protocol](#protocol) | `string` | Required | cannot be null | [alter-external-domain input](alter-external-domain-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/properties/protocol") |

## domain



`domain`

* is required

* Type: `string` ([User domain name](alter-external-domain-input-properties-user-domain-name.md))

* cannot be null

* defined in: [alter-external-domain input](alter-external-domain-input-properties-user-domain-name.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/properties/domain")

### domain Type

`string` ([User domain name](alter-external-domain-input-properties-user-domain-name.md))

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](alter-external-domain-input-properties-provider-protocol.md))

* cannot be null

* defined in: [alter-external-domain input](alter-external-domain-input-properties-provider-protocol.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/properties/protocol")

### protocol Type

`string` ([Provider protocol](alter-external-domain-input-properties-provider-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

# alter-external-domain input Definitions

## Definitions group additional-properties-of-ldap

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap"}
```

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                               |
| :------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [bind\_dn](#bind_dn)             | `string`  | Required | cannot be null | [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_dn")             |
| [bind\_password](#bind_password) | `string`  | Required | cannot be null | [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_password") |
| [tls](#tls)                      | `boolean` | Required | cannot be null | [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls")                     |
| [tls\_verify](#tls_verify)       | `boolean` | Required | cannot be null | [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls_verify")       |

### bind\_dn



`bind_dn`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_dn")

#### bind\_dn Type

`string`

#### bind\_dn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### bind\_password



`bind_password`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_password")

#### bind\_password Type

`string`

#### bind\_password Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### tls



`tls`

* is required

* Type: `boolean`

* cannot be null

* defined in: [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls")

#### tls Type

`boolean`

### tls\_verify



`tls_verify`

* is required

* Type: `boolean`

* cannot be null

* defined in: [alter-external-domain input](alter-external-domain-input-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/alter-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls_verify")

#### tls\_verify Type

`boolean`
