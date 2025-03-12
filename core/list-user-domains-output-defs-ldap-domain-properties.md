# LDAP domain properties Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap
```

Additional required properties of LDAP-based domains

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-user-domains-output.json\*](cluster/list-user-domains-output.json "open original schema") |

## additional-properties-of-ldap Type

`object` ([LDAP domain properties](list-user-domains-output-defs-ldap-domain-properties.md))

# additional-properties-of-ldap Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                      |
| :------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [schema](#schema)                | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/schema") |
| [base\_dn](#base_dn)             | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-base_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/base_dn")             |
| [bind\_dn](#bind_dn)             | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_dn")             |
| [bind\_password](#bind_password) | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_password") |
| [tls](#tls)                      | `boolean` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls")                     |
| [tls\_verify](#tls_verify)       | `boolean` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls_verify")       |

## schema



`schema`

* is required

* Type: `string` ([LDAP database schema](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/schema")

### schema Type

`string` ([LDAP database schema](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md))

### schema Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value       | Explanation |
| :---------- | :---------- |
| `"ad"`      |             |
| `"rfc2307"` |             |

## base\_dn



`base_dn`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-base_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/base_dn")

### base\_dn Type

`string`

### base\_dn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## bind\_dn



`bind_dn`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_dn")

### bind\_dn Type

`string`

### bind\_dn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## bind\_password



`bind_password`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_password")

### bind\_password Type

`string`

### bind\_password Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## tls



`tls`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls")

### tls Type

`boolean`

## tls\_verify



`tls_verify`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls_verify")

### tls\_verify Type

`boolean`
