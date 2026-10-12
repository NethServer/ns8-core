# LDAP domain properties Schema

```txt
http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap
```

Additional required properties of LDAP-based domains

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-external-domain-input.json\*](cluster/add-external-domain-input.json "open original schema") |

## additional-properties-of-ldap Type

`object` ([LDAP domain properties](add-external-domain-input-defs-ldap-domain-properties.md))

# additional-properties-of-ldap Properties

| Property                        | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                         |
| :------------------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [schema](#schema)               | `string`  | Optional | can be null    | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/schema") |
| [base_dn](#base_dn)             | `string`  | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/base_dn")             |
| [bind_dn](#bind_dn)             | `string`  | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_dn")             |
| [bind_password](#bind_password) | `string`  | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_password") |
| [tls](#tls)                     | `boolean` | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls")                     |
| [tls_verify](#tls_verify)       | `boolean` | Required | cannot be null | [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls_verify")       |

## schema

The LDAP schema is probed automatically if the value is `null` or the property is missing

`schema`

* is optional

* Type: `string` ([LDAP database schema](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md))

* can be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/schema")

### schema Type

`string` ([LDAP database schema](add-external-domain-input-defs-ldap-domain-properties-properties-ldap-database-schema.md))

### schema Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value       | Explanation |
| :---------- | :---------- |
| `"ad"`      |             |
| `"rfc2307"` |             |

## base_dn

The LDAP base DN is probed automatically if the value is `""` (empty string)

`base_dn`

* is required

* Type: `string` ([Base DN](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md))

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/base_dn")

### base_dn Type

`string` ([Base DN](add-external-domain-input-defs-ldap-domain-properties-properties-base-dn.md))

## bind_dn



`bind_dn`

* is required

* Type: `string`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_dn")

### bind_dn Type

`string`

### bind_dn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## bind_password



`bind_password`

* is required

* Type: `string`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/bind_password")

### bind_password Type

`string`

### bind_password Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## tls



`tls`

* is required

* Type: `boolean`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls")

### tls Type

`boolean`

## tls_verify



`tls_verify`

* is required

* Type: `boolean`

* cannot be null

* defined in: [add-external-domain input](add-external-domain-input-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/add-external-domain-input.json#/$defs/additional-properties-of-ldap/properties/tls_verify")

### tls_verify Type

`boolean`
