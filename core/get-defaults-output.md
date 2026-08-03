# get-defaults output Schema

```txt
http://schema.nethserver.org/cluster/get-defaults-output.json
```

Output schema of the get-defaults action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-defaults-output.json](cluster/get-defaults-output.json "open original schema") |

## get-defaults output Type

`object` ([get-defaults output](get-defaults-output.md))

## get-defaults output Examples

```json
{
  "has_additional_disks": false,
  "vpn": {
    "host": "server.nethserver.org",
    "port": 55820,
    "network": "10.5.4.0/24"
  }
}
```

# get-defaults output Properties

| Property                                        | Type      | Required | Nullable       | Defined by                                                                                                                                                                     |
| :---------------------------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [has\_additional\_disks](#has_additional_disks) | `boolean` | Optional | cannot be null | [get-defaults output](get-defaults-output-properties-has_additional_disks.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/has_additional_disks") |
| [vpn](#vpn)                                     | `object`  | Optional | cannot be null | [get-defaults output](get-defaults-output-properties-vpn.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn")                                   |

## has\_additional\_disks

True if the OS has additional disks for application data

`has_additional_disks`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [get-defaults output](get-defaults-output-properties-has_additional_disks.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/has_additional_disks")

### has\_additional\_disks Type

`boolean`

## vpn

WireGuard VPN defaults

`vpn`

* is optional

* Type: `object` ([Details](get-defaults-output-properties-vpn.md))

* cannot be null

* defined in: [get-defaults output](get-defaults-output-properties-vpn.md "http://schema.nethserver.org/cluster/get-defaults-output.json#/properties/vpn")

### vpn Type

`object` ([Details](get-defaults-output-properties-vpn.md))
