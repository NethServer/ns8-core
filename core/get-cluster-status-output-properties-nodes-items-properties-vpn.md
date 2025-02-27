# Untitled object in get-cluster-status output Schema

```txt
http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn
```

WireGuard VPN details

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-cluster-status-output.json\*](cluster/get-cluster-status-output.json "open original schema") |

## vpn Type

`object` ([Details](get-cluster-status-output-properties-nodes-items-properties-vpn.md))

# vpn Properties

| Property                     | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                                 |
| :--------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ip\_address](#ip_address)   | `string`  | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-ip_address.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/ip_address")   |
| [public\_key](#public_key)   | `string`  | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-public_key.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/public_key")   |
| [listen\_port](#listen_port) | `integer` | Optional | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-listen_port.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/listen_port") |
| [endpoint](#endpoint)        | `string`  | Optional | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-endpoint.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/endpoint")       |
| [lastseen](#lastseen)        | `integer` | Optional | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-lastseen.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/lastseen")       |
| [sent](#sent)                | `integer` | Optional | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-sent.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/sent")               |
| [rcvd](#rcvd)                | `integer` | Optional | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-rcvd.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/rcvd")               |
| [keepalive](#keepalive)      | `integer` | Optional | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-keepalive.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/keepalive")     |

## ip\_address

VPN node IP

`ip_address`

* is required

* Type: `string`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-ip_address.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/ip_address")

### ip\_address Type

`string`

## public\_key

WireGuard VPN publick key

`public_key`

* is required

* Type: `string`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-public_key.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/public_key")

### public\_key Type

`string`

## listen\_port

WireGuard Listen port number, where available on the leader node

`listen_port`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-listen_port.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/listen_port")

### listen\_port Type

`integer`

## endpoint

WireGuard VPN endpoint

`endpoint`

* is optional

* Type: `string`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-endpoint.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/endpoint")

### endpoint Type

`string`

## lastseen

Timestamp of last peer hanshake

`lastseen`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-lastseen.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/lastseen")

### lastseen Type

`integer`

## sent

Bytes sent to the peer

`sent`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-sent.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/sent")

### sent Type

`integer`

## rcvd

Bytes received from the peer

`rcvd`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-rcvd.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/rcvd")

### rcvd Type

`integer`

## keepalive

Seconds of persistent keepalive

`keepalive`

* is optional

* Type: `integer`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes-items-properties-vpn-properties-keepalive.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes/items/properties/vpn/properties/keepalive")

### keepalive Type

`integer`
