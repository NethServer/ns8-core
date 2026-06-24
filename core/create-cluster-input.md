# create-cluster input Schema

```txt
http://schema.nethserver.org/cluster/create-cluster-input.json
```

Provide basic information required by the new cluster initialization procedure

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [create-cluster-input.json](cluster/create-cluster-input.json "open original schema") |

## create-cluster input Type

`object` ([create-cluster input](create-cluster-input.md))

## create-cluster input Examples

```json
{
  "network": "10.5.4.0/24",
  "endpoint": "fc1.dp.nethserver.net:55820"
}
```

# create-cluster input Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                |
| :-------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [network](#network)   | `string` | Required | cannot be null | [create-cluster input](cluster-definitions-ipv4-cidr.md "http://schema.nethserver.org/cluster/create-cluster-input.json#/properties/network")             |
| [endpoint](#endpoint) | `string` | Required | cannot be null | [create-cluster input](create-cluster-input-properties-endpoint.md "http://schema.nethserver.org/cluster/create-cluster-input.json#/properties/endpoint") |

## network

IPv4 with netmask in CIDR notation

`network`

* is required

* Type: `string` ([IPv4 CIDR](cluster-definitions-ipv4-cidr.md))

* cannot be null

* defined in: [create-cluster input](cluster-definitions-ipv4-cidr.md "http://schema.nethserver.org/cluster/create-cluster-input.json#/properties/network")

### network Type

`string` ([IPv4 CIDR](cluster-definitions-ipv4-cidr.md))

### network Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$
```

[try pattern](https://regexr.com/?expression=%5E\(%5B0-9%5D%7B1%2C3%7D%5C.\)%7B3%7D%5B0-9%5D%7B1%2C3%7D\(%5C%2F\(%5B0-9%5D%7C%5B1-2%5D%5B0-9%5D%7C3%5B0-2%5D\)\)%3F%24 "try regular expression with regexr.com")

### network Examples

```json
"10.5.4.0/24"
```

```json
"192.168.73.0/24"
```

## endpoint

The public WireGuard VPN endpoint in `hostname:port` form. The given port may differ from 55820 as it depends on external configurations (i.e. port-forwards)

`endpoint`

* is required

* Type: `string`

* cannot be null

* defined in: [create-cluster input](create-cluster-input-properties-endpoint.md "http://schema.nethserver.org/cluster/create-cluster-input.json#/properties/endpoint")

### endpoint Type

`string`

### endpoint Examples

```json
"my.example.com:55820"
```

```json
"1.2.3.4:60000"
```
