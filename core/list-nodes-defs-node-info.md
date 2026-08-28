# Untitled object in list-nodes output Schema

```txt
http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-nodes.json\*](cluster/list-nodes.json "open original schema") |

## node-info Type

`object` ([Details](list-nodes-defs-node-info.md))

# node-info Properties

| Property                                              | Type          | Required | Nullable       | Defined by                                                                                                                                                                                      |
| :---------------------------------------------------- | :------------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id)                                  | `integer`     | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-node_id.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/node_id")                                 |
| [role](#role)                                         | Not specified | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-role.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/role")                                       |
| [ui\_name](#ui_name)                                  | `string`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-ui_name.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/ui_name")                                 |
| [vpn\_endpoint](#vpn_endpoint)                        | `string`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-vpn_endpoint.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/vpn_endpoint")                       |
| [vpn\_ip\_address](#vpn_ip_address)                   | `string`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-vpn_ip_address.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/vpn_ip_address")                   |
| [vpn\_listen\_port](#vpn_listen_port)                 | `string`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-vpn_listen_port.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/vpn_listen_port")                 |
| [app\_count](#app_count)                              | `integer`     | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-app_count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/app_count")                             |
| [os\_release](#os_release)                            | `object`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-os_release.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/os_release")                           |
| [network\_interface\_count](#network_interface_count) | `integer`     | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-network_interface_count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/network_interface_count") |
| [fqdn](#fqdn)                                         | `string`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-fqdn.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/fqdn")                                       |
| [main\_ip](#main_ip)                                  | `string`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-main_ip.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/main_ip")                                 |
| [load](#load)                                         | `object`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-load.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/load")                                       |
| [cpu](#cpu)                                           | `object`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-cpu.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu")                                         |
| [memory](#memory)                                     | `object`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-memory.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory")                                   |
| [swap](#swap)                                         | `object`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-swap.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap")                                       |
| [additional\_disk\_count](#additional_disk_count)     | `number`      | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-additional_disk_count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/additional_disk_count")     |
| [disks](#disks)                                       | `array`       | Optional | cannot be null | [list-nodes output](list-nodes-defs-node-info-properties-disks.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks")                                     |

## node\_id



`node_id`

* is optional

* Type: `integer`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-node_id.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/node_id")

### node\_id Type

`integer`

## role



`role`

* is optional

* Type: unknown

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-role.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/role")

### role Type

unknown

### role Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value            | Explanation |
| :--------------- | :---------- |
| `"ns7migration"` |             |
| `"leader"`       |             |
| `"worker"`       |             |

## ui\_name



`ui_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-ui_name.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/ui_name")

### ui\_name Type

`string`

## vpn\_endpoint



`vpn_endpoint`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-vpn_endpoint.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/vpn_endpoint")

### vpn\_endpoint Type

`string`

## vpn\_ip\_address



`vpn_ip_address`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-vpn_ip_address.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/vpn_ip_address")

### vpn\_ip\_address Type

`string`

## vpn\_listen\_port



`vpn_listen_port`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-vpn_listen_port.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/vpn_listen_port")

### vpn\_listen\_port Type

`string`

## app\_count

Number of applications installed on the node

`app_count`

* is optional

* Type: `integer`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-app_count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/app_count")

### app\_count Type

`integer`

## os\_release



`os_release`

* is optional

* Type: `object` ([Details](list-nodes-defs-node-info-properties-os_release.md))

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-os_release.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/os_release")

### os\_release Type

`object` ([Details](list-nodes-defs-node-info-properties-os_release.md))

## network\_interface\_count



`network_interface_count`

* is optional

* Type: `integer`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-network_interface_count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/network_interface_count")

### network\_interface\_count Type

`integer`

## fqdn



`fqdn`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-fqdn.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/fqdn")

### fqdn Type

`string`

## main\_ip



`main_ip`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-main_ip.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/main_ip")

### main\_ip Type

`string`

## load

Average load

`load`

* is optional

* Type: `object` ([Details](list-nodes-defs-node-info-properties-load.md))

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-load.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/load")

### load Type

`object` ([Details](list-nodes-defs-node-info-properties-load.md))

## cpu

Average load

`cpu`

* is optional

* Type: `object` ([Details](list-nodes-defs-node-info-properties-cpu.md))

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-cpu.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/cpu")

### cpu Type

`object` ([Details](list-nodes-defs-node-info-properties-cpu.md))

## memory



`memory`

* is optional

* Type: `object` ([Details](list-nodes-defs-node-info-properties-memory.md))

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-memory.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/memory")

### memory Type

`object` ([Details](list-nodes-defs-node-info-properties-memory.md))

## swap



`swap`

* is optional

* Type: `object` ([Details](list-nodes-defs-node-info-properties-swap.md))

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-swap.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/swap")

### swap Type

`object` ([Details](list-nodes-defs-node-info-properties-swap.md))

## additional\_disk\_count

Number of additional disks available on the node for application volume mapping

`additional_disk_count`

* is optional

* Type: `number`

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-additional_disk_count.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/additional_disk_count")

### additional\_disk\_count Type

`number`

## disks

Storage information by partition

`disks`

* is optional

* Type: `object[]` ([Details](list-nodes-defs-node-info-properties-disks-items.md))

* cannot be null

* defined in: [list-nodes output](list-nodes-defs-node-info-properties-disks.md "http://schema.nethserver.org/cluster/list-nodes.json#/$defs/node-info/properties/disks")

### disks Type

`object[]` ([Details](list-nodes-defs-node-info-properties-disks-items.md))
