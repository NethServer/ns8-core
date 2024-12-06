# add-node output Schema

```txt
http://schema.nethserver.org/cluster/add-node-output.json
```

Output schema of the add-node action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-node-output.json](cluster/add-node-output.json "open original schema") |

## add-node output Type

`object` ([add-node output](add-node-output.md))

## add-node output Examples

```json
{
  "node_id": 2,
  "ip_address": "10.5.4.2",
  "leader_public_key": "aoBydmHg2WWv5OqYM1ZegcCdC6dev+IHnoNv3ftNY2U=",
  "network": "10.5.4.0/24",
  "leader_ip_address": "10.5.4.1",
  "leader_endpoint": "fc1.example.com:55820"
}
```

# add-node output Properties

| Property                                  | Type      | Required | Nullable       | Defined by                                                                                                                                                              |
| :---------------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id)                      | `integer` | Required | cannot be null | [add-node output](add-node-output-properties-node-identifier.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/node_id")                        |
| [ip\_address](#ip_address)                | `string`  | Required | cannot be null | [add-node output](add-node-output-properties-node-ip-address.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/ip_address")                     |
| [network](#network)                       | `string`  | Required | cannot be null | [add-node output](cluster-definitions-ipv4-cidr.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/network")                                     |
| [leader\_public\_key](#leader_public_key) | `string`  | Required | cannot be null | [add-node output](add-node-output-properties-leader-vpn-key.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/leader_public_key")               |
| [leader\_ip\_address](#leader_ip_address) | `string`  | Required | cannot be null | [add-node output](add-node-output-properties-leader-ip-address.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/leader_ip_address")            |
| [leader\_endpoint](#leader_endpoint)      | `string`  | Required | cannot be null | [add-node output](add-node-output-properties-leader-public-vpn-udp-endpoint.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/leader_endpoint") |

## node\_id

Allocated identifier of the new node

`node_id`

* is required

* Type: `integer` ([Node identifier](add-node-output-properties-node-identifier.md))

* cannot be null

* defined in: [add-node output](add-node-output-properties-node-identifier.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/node_id")

### node\_id Type

`integer` ([Node identifier](add-node-output-properties-node-identifier.md))

## ip\_address

VPN IP address of the added node

`ip_address`

* is required

* Type: `string` ([Node IP address](add-node-output-properties-node-ip-address.md))

* cannot be null

* defined in: [add-node output](add-node-output-properties-node-ip-address.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/ip_address")

### ip\_address Type

`string` ([Node IP address](add-node-output-properties-node-ip-address.md))

### ip\_address Constraints

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")

## network

IPv4 with netmask in CIDR notation

`network`

* is required

* Type: `string` ([IPv4 CIDR](cluster-definitions-ipv4-cidr.md))

* cannot be null

* defined in: [add-node output](cluster-definitions-ipv4-cidr.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/network")

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

## leader\_public\_key

WireGuard public key of the leader node

`leader_public_key`

* is required

* Type: `string` ([Leader VPN key](add-node-output-properties-leader-vpn-key.md))

* cannot be null

* defined in: [add-node output](add-node-output-properties-leader-vpn-key.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/leader_public_key")

### leader\_public\_key Type

`string` ([Leader VPN key](add-node-output-properties-leader-vpn-key.md))

## leader\_ip\_address

VPN IP address of the leader node

`leader_ip_address`

* is required

* Type: `string` ([Leader IP address](add-node-output-properties-leader-ip-address.md))

* cannot be null

* defined in: [add-node output](add-node-output-properties-leader-ip-address.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/leader_ip_address")

### leader\_ip\_address Type

`string` ([Leader IP address](add-node-output-properties-leader-ip-address.md))

### leader\_ip\_address Constraints

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")

## leader\_endpoint

WireGuard public UDP endpoint address to contact the leader node. Syntax is `hostname:port` or `IP:port`

`leader_endpoint`

* is required

* Type: `string` ([Leader public VPN UDP endpoint](add-node-output-properties-leader-public-vpn-udp-endpoint.md))

* cannot be null

* defined in: [add-node output](add-node-output-properties-leader-public-vpn-udp-endpoint.md "http://schema.nethserver.org/cluster/add-node-output.json#/properties/leader_endpoint")

### leader\_endpoint Type

`string` ([Leader public VPN UDP endpoint](add-node-output-properties-leader-public-vpn-udp-endpoint.md))

### leader\_endpoint Examples

```json
"my.example.com:55820"
```

```json
"1.2.3.4:55820"
```
