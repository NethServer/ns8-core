# add-module input Schema

```txt
http://schema.nethserver.org/node/add-module-input.json
```

Install a module on the worker node

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-module-input.json](node/add-module-input.json "open original schema") |

## add-module input Type

`object` ([add-module input](add-module-input.md))

## add-module input Examples

```json
{
  "module_id": "mymodule2",
  "is_rootfull": false,
  "environment": {
    "IMAGE_URL": "ghcr.io/nethserver/mymodule:v2.3.2"
  }
}
```

# add-module input Properties

| Property                                | Type      | Required | Nullable       | Defined by                                                                                                                                                         |
| :-------------------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [environment](#environment)             | `object`  | Required | cannot be null | [add-module input](add-module-input-properties-initial-module-environment.md "http://schema.nethserver.org/node/add-module-input.json#/properties/environment")    |
| [is\_rootfull](#is_rootfull)            | `boolean` | Required | cannot be null | [add-module input](add-module-input-properties-required-type-of-installation.md "http://schema.nethserver.org/node/add-module-input.json#/properties/is_rootfull") |
| [module\_id](#module_id)                | `string`  | Required | cannot be null | [add-module input](add-module-input-properties-module-identifier.md "http://schema.nethserver.org/node/add-module-input.json#/properties/module_id")               |
| [tcp\_ports\_demand](#tcp_ports_demand) | `number`  | Required | cannot be null | [add-module input](add-module-input-properties-tcp-ports.md "http://schema.nethserver.org/node/add-module-input.json#/properties/tcp_ports_demand")                |
| [udp\_ports\_demand](#udp_ports_demand) | `number`  | Required | cannot be null | [add-module input](add-module-input-properties-udp-ports.md "http://schema.nethserver.org/node/add-module-input.json#/properties/udp_ports_demand")                |

## environment

Assign initial values to the module environment

`environment`

* is required

* Type: `object` ([Initial module environment](add-module-input-properties-initial-module-environment.md))

* cannot be null

* defined in: [add-module input](add-module-input-properties-initial-module-environment.md "http://schema.nethserver.org/node/add-module-input.json#/properties/environment")

### environment Type

`object` ([Initial module environment](add-module-input-properties-initial-module-environment.md))

## is\_rootfull

If `true` the module is installed in rootfull mode.
Its agent runs as root and its containers can be granted **privileged access**.

`is_rootfull`

* is required

* Type: `boolean` ([Required type of installation](add-module-input-properties-required-type-of-installation.md))

* cannot be null

* defined in: [add-module input](add-module-input-properties-required-type-of-installation.md "http://schema.nethserver.org/node/add-module-input.json#/properties/is_rootfull")

### is\_rootfull Type

`boolean` ([Required type of installation](add-module-input-properties-required-type-of-installation.md))

## module\_id



`module_id`

* is required

* Type: `string` ([Module identifier](add-module-input-properties-module-identifier.md))

* cannot be null

* defined in: [add-module input](add-module-input-properties-module-identifier.md "http://schema.nethserver.org/node/add-module-input.json#/properties/module_id")

### module\_id Type

`string` ([Module identifier](add-module-input-properties-module-identifier.md))

### module\_id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### module\_id Examples

```json
"dokuwiki1"
```

```json
"traefik3"
```

```json
"samba1"
```

## tcp\_ports\_demand

Number of TCP ports that will be allocate

`tcp_ports_demand`

* is required

* Type: `number` ([TCP ports](add-module-input-properties-tcp-ports.md))

* cannot be null

* defined in: [add-module input](add-module-input-properties-tcp-ports.md "http://schema.nethserver.org/node/add-module-input.json#/properties/tcp_ports_demand")

### tcp\_ports\_demand Type

`number` ([TCP ports](add-module-input-properties-tcp-ports.md))

## udp\_ports\_demand

Number of UDP ports that will be allocate

`udp_ports_demand`

* is required

* Type: `number` ([UDP ports](add-module-input-properties-udp-ports.md))

* cannot be null

* defined in: [add-module input](add-module-input-properties-udp-ports.md "http://schema.nethserver.org/node/add-module-input.json#/properties/udp_ports_demand")

### udp\_ports\_demand Type

`number` ([UDP ports](add-module-input-properties-udp-ports.md))
