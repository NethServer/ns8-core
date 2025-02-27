# Untitled object in get-defaults output Schema

```txt
http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn
```

WireGuard VPN defaults

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-defaults-output.json\*](cluster/get-defaults-output.json "open original schema") |

## vpn Type

`object` ([Details](get-defaults-output-properties-vpn.md))

# vpn Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                                                         |
| :------------------ | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [host](#host)       | `string`  | Required | cannot be null | [get-defaults output](get-defaults-output-properties-vpn-properties-host.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn/properties/host")       |
| [port](#port)       | `integer` | Required | cannot be null | [get-defaults output](get-defaults-output-properties-vpn-properties-port.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn/properties/port")       |
| [network](#network) | `string`  | Required | cannot be null | [get-defaults output](get-defaults-output-properties-vpn-properties-network.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn/properties/network") |

## host

WireGuard public endpoint name

`host`

* is required

* Type: `string`

* cannot be null

* defined in: [get-defaults output](get-defaults-output-properties-vpn-properties-host.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn/properties/host")

### host Type

`string`

## port

WireGuard listen port

`port`

* is required

* Type: `integer`

* cannot be null

* defined in: [get-defaults output](get-defaults-output-properties-vpn-properties-port.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn/properties/port")

### port Type

`integer`

## network

Wireguard peer network

`network`

* is required

* Type: `string`

* cannot be null

* defined in: [get-defaults output](get-defaults-output-properties-vpn-properties-network.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn/properties/network")

### network Type

`string`
