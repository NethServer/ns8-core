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

| Property                  | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                         |
| :------------------------ | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ip_address](#ip_address) | `string`  | Required | cannot be null | [update-routes input](update-routes-input-definitions-changelist-items-properties-ip-address.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/ip_address")               |
| [node_id](#node_id)       | `integer` | Required | cannot be null | [update-routes input](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/node_id") |

## ip_address

IP address to add or remove. It should be local to the node.

`ip_address`

* is required

* Type: `string` ([IP address](update-routes-input-definitions-changelist-items-properties-ip-address.md))

* cannot be null

* defined in: [update-routes input](update-routes-input-definitions-changelist-items-properties-ip-address.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/ip_address")

### ip_address Type

`string` ([IP address](update-routes-input-definitions-changelist-items-properties-ip-address.md))

### ip_address Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")

## node_id

Node ID used as route next-hop

`node_id`

* is required

* Type: `integer` ([Destination node identifier](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md))

* cannot be null

* defined in: [update-routes input](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md "http://schema.nethserver.org/cluster/update-routes-input.json#/definitions/changeList/items/properties/node_id")

### node_id Type

`integer` ([Destination node identifier](update-routes-input-definitions-changelist-items-properties-destination-node-identifier.md))

### node_id Constraints

**minimum (exclusive)**: the value of this number must be greater than: `0`
