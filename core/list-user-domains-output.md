# list-user-domains output Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json
```

Quickly get the user domains list and their basic configuration

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-user-domains-output.json](cluster/list-user-domains-output.json "open original schema") |

## list-user-domains output Type

`object` ([list-user-domains output](list-user-domains-output.md))

## list-user-domains output Examples

```json
{
  "unconfigured_domains": [
    {
      "module_id": "samba1",
      "image_name": "samba",
      "image_url": "ghcr.io/nethserver/samba:latest"
    }
  ],
  "domains": [
    {
      "name": "sandbox.example",
      "location": "internal",
      "protocol": "ldap",
      "schema": "rfc2307",
      "base_dn": "dc=sandbox,dc=example",
      "bind_dn": "cn=ldapservice,dc=sandbox,dc=example",
      "bind_password": "S3cr3t!",
      "tls": false,
      "tls_verify": false,
      "counters": {
        "users": null,
        "groups": null
      },
      "providers": [
        {
          "id": "openldap1",
          "ui_name": "",
          "node": 1,
          "file_server": false,
          "host": "10.110.32.2",
          "port": 20003
        },
        {
          "id": "openldap2",
          "ui_name": "",
          "node": 2,
          "file_server": false,
          "host": "10.110.32.3",
          "port": 20002
        }
      ]
    },
    {
      "name": "company.org",
      "location": "external",
      "protocol": "ldap",
      "schema": "rfc2307",
      "base_dn": "dc=company,dc=org",
      "bind_dn": "cn=ns8cluster,dc=company,dc=org",
      "bind_password": "OtherS3cr3t!",
      "tls": true,
      "tls_verify": true,
      "counters": {
        "users": 15,
        "groups": 3
      },
      "providers": [
        {
          "id": "ldap-primary.company.org",
          "ui_name": "Company LDAP primary",
          "node": null,
          "file_server": false,
          "host": "ldap-master.company.org",
          "port": 636
        },
        {
          "id": "ldap-replica.company.org",
          "ui_name": "Company LDAP replica",
          "node": null,
          "file_server": false,
          "host": "ldap-replica.company.org",
          "port": 636
        }
      ]
    }
  ]
}
```

# list-user-domains output Properties

| Property                                       | Type    | Required | Nullable       | Defined by                                                                                                                                                                                    |
| :--------------------------------------------- | :------ | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [unconfigured\_domains](#unconfigured_domains) | `array` | Required | cannot be null | [list-user-domains output](list-user-domains-output-properties-unconfigured_domains.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/properties/unconfigured_domains") |
| [domains](#domains)                            | `array` | Required | cannot be null | [list-user-domains output](list-user-domains-output-properties-domains.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/properties/domains")                           |

## unconfigured\_domains



`unconfigured_domains`

* is required

* Type: `object[]` ([Unconfigured domain](list-user-domains-output-defs-unconfigured-domain.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-properties-unconfigured_domains.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/properties/unconfigured_domains")

### unconfigured\_domains Type

`object[]` ([Unconfigured domain](list-user-domains-output-defs-unconfigured-domain.md))

## domains



`domains`

* is required

* Type: `object[]` ([User domain](list-user-domains-output-defs-user-domain.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-properties-domains.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/properties/domains")

### domains Type

`object[]` ([User domain](list-user-domains-output-defs-user-domain.md))

### domains Constraints

**minimum number of items**: the minimum number of items for this array is: `0`

# list-user-domains output Definitions

## Definitions group user-domain

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain"}
```

| Property                         | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                              |
| :------------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)                    | `string` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-name.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/name")                        |
| [location](#location)            | `string` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/location") |
| [counters](#counters)            | `object` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters")                |
| [hidden\_users](#hidden_users)   | `array`  | Optional | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_users.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_users")        |
| [hidden\_groups](#hidden_groups) | `array`  | Optional | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_groups.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_groups")      |
| [protocol](#protocol)            | `string` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/protocol")       |
| [providers](#providers)          | `array`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-account-providers.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/providers")      |

### name



`name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-name.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/name")

#### name Type

`string`

#### name Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### location

Type of domain hosting. Set to `internal` if the domain is hosted by the cluster, `external` if the domain is provided by a remote service

`location`

* is required

* Type: `string` ([Domain hosting location](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/location")

#### location Type

`string` ([Domain hosting location](list-user-domains-output-defs-user-domain-properties-domain-hosting-location.md))

#### location Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value        | Explanation |
| :----------- | :---------- |
| `"internal"` |             |
| `"external"` |             |

### counters

The cached number of users and groups returned by their respective last API calls

`counters`

* is required

* Type: `object` ([Counters](list-user-domains-output-defs-user-domain-properties-counters.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters")

#### counters Type

`object` ([Counters](list-user-domains-output-defs-user-domain-properties-counters.md))

### hidden\_users

List of users that are not visible from UI and from applications

`hidden_users`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_users.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_users")

#### hidden\_users Type

`string[]`

### hidden\_groups

List of groups that are not visible from UI and from applications

`hidden_groups`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-hidden_groups.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/hidden_groups")

#### hidden\_groups Type

`string[]`

### protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](list-user-domains-output-defs-user-domain-properties-provider-protocol.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/protocol")

#### protocol Type

`string` ([Provider protocol](list-user-domains-output-defs-user-domain-properties-provider-protocol.md))

#### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

### providers

Backend system and replicas providing the services of the user domain

`providers`

* is required

* Type: an array of merged types ([Details](list-user-domains-output-defs-user-domain-properties-account-providers-items.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-account-providers.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/providers")

#### providers Type

an array of merged types ([Details](list-user-domains-output-defs-user-domain-properties-account-providers-items.md))

#### providers Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

## Definitions group unconfigured-domain

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain"}
```

| Property                 | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                              |
| :----------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id) | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/module_id")      |
| [node](#node)            | `integer` | Required | can be null    | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/node")             |
| [location](#location-1)  | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/location") |
| [protocol](#protocol-1)  | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/protocol")       |
| [schema](#schema)        | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/schema")      |

### module\_id

e.g. `samba1`

`module_id`

* is required

* Type: `string` ([Module identifier](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/module_id")

#### module\_id Type

`string` ([Module identifier](list-user-domains-output-defs-unconfigured-domain-properties-module-identifier.md))

### node

The node number, e.g. `1`

`node`

* is required

* Type: `integer` ([Node identifier](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md))

* can be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/node")

#### node Type

`integer` ([Node identifier](list-user-domains-output-defs-unconfigured-domain-properties-node-identifier.md))

#### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`

### location

Type of domain hosting. Set to `internal` if the domain is hosted by the cluster, `external` if the domain is provided by a remote service

`location`

* is required

* Type: `string` ([Domain hosting location](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/location")

#### location Type

`string` ([Domain hosting location](list-user-domains-output-defs-unconfigured-domain-properties-domain-hosting-location.md))

#### location Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value        | Explanation |
| :----------- | :---------- |
| `"internal"` |             |
| `"external"` |             |

### protocol

Protocol used to communicate with the domain providers.

`protocol`

* is required

* Type: `string` ([Provider protocol](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/protocol")

#### protocol Type

`string` ([Provider protocol](list-user-domains-output-defs-unconfigured-domain-properties-provider-protocol.md))

#### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value    | Explanation |
| :------- | :---------- |
| `"ldap"` |             |

### schema



`schema`

* is required

* Type: `string` ([LDAP database schema](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/unconfigured-domain/properties/schema")

#### schema Type

`string` ([LDAP database schema](list-user-domains-output-defs-unconfigured-domain-properties-ldap-database-schema.md))

#### schema Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value       | Explanation |
| :---------- | :---------- |
| `"ad"`      |             |
| `"rfc2307"` |             |

## Definitions group ldap-provider

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider"}
```

| Property                     | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                                     |
| :--------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)                    | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-id.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/id")                                                       |
| [ui\_name](#ui_name)         | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-ui_name.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/ui_name")                                             |
| [node](#node-1)              | `integer` | Required | can be null    | [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-node.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/node")                                                   |
| [file\_server](#file_server) | `boolean` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-the-provider-can-be-used-as-smb-file-server-too.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/file_server") |
| [host](#host)                | `string`  | Required | can be null    | [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-host-name-or-ip-address.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/host")                                |
| [port](#port)                | `integer` | Required | can be null    | [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-tcp-port-number.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/port")                                        |

### id



`id`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-id.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/id")

#### id Type

`string`

#### id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### ui\_name



`ui_name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-ui_name.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/ui_name")

#### ui\_name Type

`string`

### node



`node`

* is required

* Type: `integer`

* can be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-node.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/node")

#### node Type

`integer`

#### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`

### file\_server



`file_server`

* is required

* Type: `boolean` ([The provider can be used as SMB file server too](list-user-domains-output-defs-ldap-account-provider-properties-the-provider-can-be-used-as-smb-file-server-too.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-the-provider-can-be-used-as-smb-file-server-too.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/file_server")

#### file\_server Type

`boolean` ([The provider can be used as SMB file server too](list-user-domains-output-defs-ldap-account-provider-properties-the-provider-can-be-used-as-smb-file-server-too.md))

### host

Can be `null` if the provider is not configured properly

`host`

* is required

* Type: `string` ([Host name or IP address](list-user-domains-output-defs-ldap-account-provider-properties-host-name-or-ip-address.md))

* can be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-host-name-or-ip-address.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/host")

#### host Type

`string` ([Host name or IP address](list-user-domains-output-defs-ldap-account-provider-properties-host-name-or-ip-address.md))

#### host Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### port

Can be `null` if the provider is not configured properly

`port`

* is required

* Type: `integer` ([TCP port number](list-user-domains-output-defs-ldap-account-provider-properties-tcp-port-number.md))

* can be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-account-provider-properties-tcp-port-number.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/ldap-provider/properties/port")

#### port Type

`integer` ([TCP port number](list-user-domains-output-defs-ldap-account-provider-properties-tcp-port-number.md))

#### port Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## Definitions group additional-properties-of-ldap

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap"}
```

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                      |
| :------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [schema](#schema-1)              | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/schema") |
| [base\_dn](#base_dn)             | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-base_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/base_dn")             |
| [bind\_dn](#bind_dn)             | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_dn")             |
| [bind\_password](#bind_password) | `string`  | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_password") |
| [tls](#tls)                      | `boolean` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls")                     |
| [tls\_verify](#tls_verify)       | `boolean` | Required | cannot be null | [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls_verify")       |

### schema



`schema`

* is required

* Type: `string` ([LDAP database schema](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md))

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/schema")

#### schema Type

`string` ([LDAP database schema](list-user-domains-output-defs-ldap-domain-properties-properties-ldap-database-schema.md))

#### schema Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value       | Explanation |
| :---------- | :---------- |
| `"ad"`      |             |
| `"rfc2307"` |             |

### base\_dn



`base_dn`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-base_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/base_dn")

#### base\_dn Type

`string`

#### base\_dn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### bind\_dn



`bind_dn`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_dn.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_dn")

#### bind\_dn Type

`string`

#### bind\_dn Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### bind\_password



`bind_password`

* is required

* Type: `string`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-bind_password.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/bind_password")

#### bind\_password Type

`string`

#### bind\_password Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### tls



`tls`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls")

#### tls Type

`boolean`

### tls\_verify



`tls_verify`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-user-domains output](list-user-domains-output-defs-ldap-domain-properties-properties-tls_verify.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/additional-properties-of-ldap/properties/tls_verify")

#### tls\_verify Type

`boolean`
