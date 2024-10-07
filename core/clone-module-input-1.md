# clone-module input Schema

```txt
http://schema.nethserver.org/cluster/clone-module-input.json
```

Input schema of the clone-module action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [clone-module-input.json](cluster/clone-module-input.json "open original schema") |

## clone-module input Type

`object` ([clone-module input](clone-module-input-1.md))

## clone-module input Examples

```json
{
  "module": "dokuwiki1",
  "node": 1,
  "replace": false
}
```

# clone-module input Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                               |
| :------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module](#module)   | `string`  | Required | cannot be null | [clone-module input](clone-module-input-1-properties-module-id.md "http://schema.nethserver.org/cluster/clone-module-input.json#/properties/module")     |
| [node](#node)       | `integer` | Required | cannot be null | [clone-module input](clone-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/clone-module-input.json#/properties/node")         |
| [replace](#replace) | `boolean` | Required | cannot be null | [clone-module input](clone-module-input-1-properties-replace-flag.md "http://schema.nethserver.org/cluster/clone-module-input.json#/properties/replace") |

## module

The identifier of the module to be cloned

`module`

* is required

* Type: `string` ([Module ID](clone-module-input-1-properties-module-id.md))

* cannot be null

* defined in: [clone-module input](clone-module-input-1-properties-module-id.md "http://schema.nethserver.org/cluster/clone-module-input.json#/properties/module")

### module Type

`string` ([Module ID](clone-module-input-1-properties-module-id.md))

### module Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## node

Node identifier where the clone is created

`node`

* is required

* Type: `integer` ([Node ID](clone-module-input-1-properties-node-id.md))

* cannot be null

* defined in: [clone-module input](clone-module-input-1-properties-node-id.md "http://schema.nethserver.org/cluster/clone-module-input.json#/properties/node")

### node Type

`integer` ([Node ID](clone-module-input-1-properties-node-id.md))

### node Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## replace

If set to true, the clone receives the original UUID and the original module is erased. If false, the clone is just a copy of the original module.

`replace`

* is required

* Type: `boolean` ([Replace flag](clone-module-input-1-properties-replace-flag.md))

* cannot be null

* defined in: [clone-module input](clone-module-input-1-properties-replace-flag.md "http://schema.nethserver.org/cluster/clone-module-input.json#/properties/replace")

### replace Type

`boolean` ([Replace flag](clone-module-input-1-properties-replace-flag.md))
