# promote-node input Schema

```txt
http://schema.nethserver.org/cluster/promote-node-input.json
```

Promote a node to cluster leader

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [promote-node-input.json](cluster/promote-node-input.json "open original schema") |

## promote-node input Type

`object` ([promote-node input](promote-node-input.md))

## promote-node input Examples

```json
{
  "node_id": 3,
  "endpoint_address": "4.5.6.7",
  "endpoint_port": 55820,
  "endpoint_validation": true
}
```

# promote-node input Properties

| Property                                     | Type      | Required | Nullable       | Defined by                                                                                                                                                                     |
| :------------------------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id)                         | `integer` | Required | cannot be null | [promote-node input](promote-node-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/node_id")                      |
| [endpoint\_address](#endpoint_address)       | `string`  | Required | cannot be null | [promote-node input](promote-node-input-properties-endpoint-address.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_address")            |
| [endpoint\_port](#endpoint_port)             | `integer` | Required | cannot be null | [promote-node input](promote-node-input-properties-endpoint-udp-port-for-vpn.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_port")      |
| [endpoint\_validation](#endpoint_validation) | `boolean` | Optional | cannot be null | [promote-node input](promote-node-input-properties-endpoint-validation-flag.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_validation") |

## node\_id

The node ID of the new leader node

`node_id`

* is required

* Type: `integer` ([Node identifier](promote-node-input-properties-node-identifier.md))

* cannot be null

* defined in: [promote-node input](promote-node-input-properties-node-identifier.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/node_id")

### node\_id Type

`integer` ([Node identifier](promote-node-input-properties-node-identifier.md))

### node\_id Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## endpoint\_address

Host name or IP address where the new leader node can be reached by other nodes

`endpoint_address`

* is required

* Type: `string` ([Endpoint address](promote-node-input-properties-endpoint-address.md))

* cannot be null

* defined in: [promote-node input](promote-node-input-properties-endpoint-address.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_address")

### endpoint\_address Type

`string` ([Endpoint address](promote-node-input-properties-endpoint-address.md))

### endpoint\_address Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")

## endpoint\_port

The UDP port number where Wireguard is listening. It must be reachable by other cluster nodes

`endpoint_port`

* is required

* Type: `integer` ([Endpoint UDP port for VPN](promote-node-input-properties-endpoint-udp-port-for-vpn.md))

* cannot be null

* defined in: [promote-node input](promote-node-input-properties-endpoint-udp-port-for-vpn.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_port")

### endpoint\_port Type

`integer` ([Endpoint UDP port for VPN](promote-node-input-properties-endpoint-udp-port-for-vpn.md))

### endpoint\_port Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## endpoint\_validation

If set to "true", validate the endpoint with a connection attempt

`endpoint_validation`

* is optional

* Type: `boolean` ([Endpoint validation flag](promote-node-input-properties-endpoint-validation-flag.md))

* cannot be null

* defined in: [promote-node input](promote-node-input-properties-endpoint-validation-flag.md "http://schema.nethserver.org/cluster/promote-node-input.json#/properties/endpoint_validation")

### endpoint\_validation Type

`boolean` ([Endpoint validation flag](promote-node-input-properties-endpoint-validation-flag.md))
