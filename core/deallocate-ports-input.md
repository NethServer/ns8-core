# deallocate-ports input Schema

```txt
http://schema.nethserver.org/node/deallocate-ports-input.json
```

Deallocate TCP or UDP ports on a node

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [deallocate-ports-input.json](node/deallocate-ports-input.json "open original schema") |

## deallocate-ports input Type

`object` ([deallocate-ports input](deallocate-ports-input.md))

# deallocate-ports input Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                             |
| :----------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id) | `string` | Required | cannot be null | [deallocate-ports input](deallocate-ports-input-properties-module-identifier.md "http://schema.nethserver.org/node/deallocate-ports-input.json#/properties/module_id") |
| [protocol](#protocol)    | `string` | Required | cannot be null | [deallocate-ports input](deallocate-ports-input-properties-ports-protocol.md "http://schema.nethserver.org/node/deallocate-ports-input.json#/properties/protocol")     |

## module\_id

Ports are deallocated from the given module.

`module_id`

* is required

* Type: `string` ([Module identifier](deallocate-ports-input-properties-module-identifier.md))

* cannot be null

* defined in: [deallocate-ports input](deallocate-ports-input-properties-module-identifier.md "http://schema.nethserver.org/node/deallocate-ports-input.json#/properties/module_id")

### module\_id Type

`string` ([Module identifier](deallocate-ports-input-properties-module-identifier.md))

## protocol



`protocol`

* is required

* Type: `string` ([Ports protocol](deallocate-ports-input-properties-ports-protocol.md))

* cannot be null

* defined in: [deallocate-ports input](deallocate-ports-input-properties-ports-protocol.md "http://schema.nethserver.org/node/deallocate-ports-input.json#/properties/protocol")

### protocol Type

`string` ([Ports protocol](deallocate-ports-input-properties-ports-protocol.md))

### protocol Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value   | Explanation |
| :------ | :---------- |
| `"tcp"` |             |
| `"udp"` |             |
