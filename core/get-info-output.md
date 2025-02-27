# get-info output Schema

```txt
http://schema.nethserver.org/node/get-info-output.json
```

Output schema of the get-info action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-info-output.json](node/get-info-output.json "open original schema") |

## get-info output Type

`object` ([get-info output](get-info-output.md))

## get-info output Examples

```json
{
  "hostname": "node1.example.org",
  "vpn_port": 55820
}
```

# get-info output Properties

| Property               | Type      | Required | Nullable       | Defined by                                                                                                                              |
| :--------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------- |
| [hostname](#hostname)  | `string`  | Required | cannot be null | [get-info output](get-info-output-properties-hostname.md "http://schema.nethserver.org/node/get-info-output.json#/properties/hostname") |
| [vpn\_port](#vpn_port) | `integer` | Required | cannot be null | [get-info output](get-info-output-properties-vpn_port.md "http://schema.nethserver.org/node/get-info-output.json#/properties/vpn_port") |

## hostname

The node fully qualified hostname

`hostname`

* is required

* Type: `string`

* cannot be null

* defined in: [get-info output](get-info-output-properties-hostname.md "http://schema.nethserver.org/node/get-info-output.json#/properties/hostname")

### hostname Type

`string`

## vpn\_port

The listening port number of wg0 Wireguard interface, detected at runtime

`vpn_port`

* is required

* Type: `integer`

* cannot be null

* defined in: [get-info output](get-info-output-properties-vpn_port.md "http://schema.nethserver.org/node/get-info-output.json#/properties/vpn_port")

### vpn\_port Type

`integer`
