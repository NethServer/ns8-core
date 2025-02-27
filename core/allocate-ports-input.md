# allocate-ports input Schema

```txt
http://schema.nethserver.org/node/allocate-ports-input.json
```

Allocate TCP or UDP ports on a node

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                         |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [allocate-ports-input.json](node/allocate-ports-input.json "open original schema") |

## allocate-ports input Type

`object` ([allocate-ports input](allocate-ports-input.md))

# allocate-ports input Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                             |
| :------------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ports](#ports)                  | `integer` | Required | cannot be null | [allocate-ports input](allocate-ports-input-properties-number-of-tcpudp-ports.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/ports")      |
| [module\_id](#module_id)         | `string`  | Required | cannot be null | [allocate-ports input](allocate-ports-input-properties-module-identifier.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/module_id")       |
| [protocol](#protocol)            | `string`  | Required | cannot be null | [allocate-ports input](allocate-ports-input-properties-ports-protocol.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/protocol")           |
| [keep\_existing](#keep_existing) | `boolean` | Required | cannot be null | [allocate-ports input](allocate-ports-input-properties-keep-existing-ports.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/keep_existing") |

## ports

How many ports will be allocated on a specific node

`ports`

* is required

* Type: `integer` ([Number of TCP/UDP ports](allocate-ports-input-properties-number-of-tcpudp-ports.md))

* cannot be null

* defined in: [allocate-ports input](allocate-ports-input-properties-number-of-tcpudp-ports.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/ports")

### ports Type

`integer` ([Number of TCP/UDP ports](allocate-ports-input-properties-number-of-tcpudp-ports.md))

## module\_id

Ports are allocated to the given module.

`module_id`

* is required

* Type: `string` ([Module identifier](allocate-ports-input-properties-module-identifier.md))

* cannot be null

* defined in: [allocate-ports input](allocate-ports-input-properties-module-identifier.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/module_id")

### module\_id Type

`string` ([Module identifier](allocate-ports-input-properties-module-identifier.md))

## protocol



`protocol`

* is required

* Type: `string` ([Ports protocol](allocate-ports-input-properties-ports-protocol.md))

* cannot be null

* defined in: [allocate-ports input](allocate-ports-input-properties-ports-protocol.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/protocol")

### protocol Type

`string` ([Ports protocol](allocate-ports-input-properties-ports-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value   | Explanation |
| :------ | :---------- |
| `"tcp"` |             |
| `"udp"` |             |

## keep\_existing

If false, remove remove existing ports before allocating a new range

`keep_existing`

* is required

* Type: `boolean` ([Keep existing ports](allocate-ports-input-properties-keep-existing-ports.md))

* cannot be null

* defined in: [allocate-ports input](allocate-ports-input-properties-keep-existing-ports.md "http://schema.nethserver.org/node/allocate-ports-input.json#/properties/keep_existing")

### keep\_existing Type

`boolean` ([Keep existing ports](allocate-ports-input-properties-keep-existing-ports.md))
