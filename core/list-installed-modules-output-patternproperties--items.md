# Untitled object in Output for list-installed-modules Schema

```txt
http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-installed-modules-output.json\*](cluster/list-installed-modules-output.json "open original schema") |

## items Type

`object` ([Details](list-installed-modules-output-patternproperties--items.md))

# items Properties

| Property            | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                 |
| :------------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)           | `string` | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-id.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/id")           |
| [node](#node)       | `string` | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-node.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/node")       |
| [digest](#digest)   | `string` | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-digest.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/digest")   |
| [source](#source)   | `string` | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-source.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/source")   |
| [version](#version) | `string` | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-version.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/version") |
| [logo](#logo)       | `string` | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-logo.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/logo")       |
| [module](#module)   | `string` | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-module.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/module")   |
| [flags](#flags)     | `array`  | Required | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-flags.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/flags")     |

## id

Unique name of a module instance

`id`

* is required

* Type: `string`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-id.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/id")

### id Type

`string`

## node

Id of the node where the instance is running

`node`

* is required

* Type: `string`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-node.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/node")

### node Type

`string`

## digest

Image digest

`digest`

* is required

* Type: `string`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-digest.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/digest")

### digest Type

`string`

## source

The URL of the container image registry

`source`

* is required

* Type: `string`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-source.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/source")

### source Type

`string`

## version

A valid semantic version extracted from image tag

`version`

* is required

* Type: `string`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-version.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/version")

### version Type

`string`

## logo

The logo URL from repository cache, it can be an empty string

`logo`

* is required

* Type: `string`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-logo.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/logo")

### logo Type

`string`

## module

Name of the module of the instance

`module`

* is required

* Type: `string`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-module.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/module")

### module Type

`string`

## flags

List of flags from org.nethserver.flags image label

`flags`

* is required

* Type: `string[]`

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties--items-properties-flags.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*/items/properties/flags")

### flags Type

`string[]`
