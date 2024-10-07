# add-node input Schema

```txt
http://schema.nethserver.org/cluster/add-node-input.json
```

Input schema of the add-node action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-node-input.json](cluster/add-node-input.json "open original schema") |

## add-node input Type

`object` ([add-node input](add-node-input.md))

## add-node input Examples

```json
{
  "endpoint": "",
  "node_pwh": "b12f4f5137c4c6b4ef2efd936b88119c4af8356d4e9c3a31a18cd4a88effc41d",
  "public_key": "XCVh3GCRqMcmIgpHhsbCL7Usm7Lm7Hb7bsS5LTSz/Dc="
}
```

```json
{
  "endpoint": "",
  "node_pwh": "b12f4f5137c4c6b4ef2efd936b88119c4af8356d4e9c3a31a18cd4a88effc41d",
  "public_key": "XCVh3GCRqMcmIgpHhsbCL7Usm7Lm7Hb7bsS5LTSz/Dc=",
  "flags": [
    "nomodules"
  ]
}
```

# add-node input Properties

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                  |
| :------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------ |
| [node\_pwh](#node_pwh)     | `string` | Required | cannot be null | [add-node input](cluster-definitions-redis-password.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/node_pwh")     |
| [public\_key](#public_key) | `string` | Required | cannot be null | [add-node input](add-node-input-properties-public_key.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/public_key") |
| [endpoint](#endpoint)      | `string` | Required | cannot be null | [add-node input](add-node-input-properties-endpoint.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/endpoint")     |
| [flags](#flags)            | `array`  | Optional | cannot be null | [add-node input](add-node-input-properties-flags.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/flags")           |

## node\_pwh

Redis SHA256 hashed password for Users, Cluster, Nodes and Modules

`node_pwh`

* is required

* Type: `string` ([Redis password](cluster-definitions-redis-password.md))

* cannot be null

* defined in: [add-node input](cluster-definitions-redis-password.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/node_pwh")

### node\_pwh Type

`string` ([Redis password](cluster-definitions-redis-password.md))

### node\_pwh Constraints

**maximum length**: the maximum number of characters for this string is: `64`

**minimum length**: the minimum number of characters for this string is: `64`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[a-f0-9]{64}$
```

[try pattern](https://regexr.com/?expression=%5E%5Ba-f0-9%5D%7B64%7D%24 "try regular expression with regexr.com")

### node\_pwh Examples

```json
"73cb3858a687a8494ca3323053016282f3dad39d42cf62ca4e79dda2aac7d9ac"
```

## public\_key

Node public key of the WireGuard VPN

`public_key`

* is required

* Type: `string`

* cannot be null

* defined in: [add-node input](add-node-input-properties-public_key.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/public_key")

### public\_key Type

`string`

## endpoint

Announce the public WireGuard VPN endpoint of the node

`endpoint`

* is required

* Type: `string`

* cannot be null

* defined in: [add-node input](add-node-input-properties-endpoint.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/endpoint")

### endpoint Type

`string`

## flags

List of node flags

`flags`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [add-node input](add-node-input-properties-flags.md "http://schema.nethserver.org/cluster/add-node-input.json#/properties/flags")

### flags Type

`string[]`
