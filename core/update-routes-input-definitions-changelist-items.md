# Untitled object in update-routes input Schema

```txt
http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [update-routes-input.json\*](cluster/update-routes-input.json "open original schema") |

## items Type

`object` ([Details](update-routes-input-definitions-changelist-items.md))

# items Properties

| Property                   | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                         |
| :------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ip\_address](#ip_address) | `string`  | Required | cannot be null | [update-routes input](update-routes-input-definitions-changelist-items-properties-ip-address.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/ip_address")               |
| [node\_id](#node_id)       | `integer` | Required | cannot be null | [update-routes input](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/node_id") |

## ip\_address

IP address to add or remove. It should be local to the node.

`ip_address`

* is required

* Type: `string` ([IP address](update-routes-input-definitions-changelist-items-properties-ip-address.md))

* cannot be null

* defined in: [update-routes input](update-routes-input-definitions-changelist-items-properties-ip-address.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/ip_address")

### ip\_address Type

`string` ([IP address](update-routes-input-definitions-changelist-items-properties-ip-address.md))

### ip\_address Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")

## node\_id

Node ID used as route next-hop

`node_id`

* is required

* Type: `integer` ([Destination node identifier](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md))

* cannot be null

* defined in: [update-routes input](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/node_id")

### node\_id Type

`integer` ([Destination node identifier](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md))

### node\_id Constraints

**minimum (exclusive)**: the value of this number must be greater than: `0`
