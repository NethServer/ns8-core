# Service discovery information Schema

```txt
http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-service-providers-output.json\*](agent/list-service-providers-output.json "open original schema") |

## itmes Type

`object` ([Service discovery information](list-service-providers-output-service-discovery-information.md))

# itmes Properties

| Property                     | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                        |
| :--------------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id)     | `string` | Required | cannot be null | [list-service-providers output](list-service-providers-output-service-discovery-information-properties-module-id.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/module_id")          |
| [module\_uuid](#module_uuid) | `string` | Required | cannot be null | [list-service-providers output](list-service-providers-output-service-discovery-information-properties-module-uuid.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/module_uuid")      |
| [node](#node)                | `string` | Required | cannot be null | [list-service-providers output](list-service-providers-output-service-discovery-information-properties-node-id.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/node")                 |
| [transport](#transport)      | `string` | Required | cannot be null | [list-service-providers output](list-service-providers-output-service-discovery-information-properties-transport-protocol.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/transport") |
| [ui\_name](#ui_name)         | `string` | Required | can be null    | [list-service-providers output](list-service-providers-output-service-discovery-information-properties-ui-name.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/ui_name")              |

## module\_id



`module_id`

* is required

* Type: `string` ([Module ID](list-service-providers-output-service-discovery-information-properties-module-id.md))

* cannot be null

* defined in: [list-service-providers output](list-service-providers-output-service-discovery-information-properties-module-id.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/module_id")

### module\_id Type

`string` ([Module ID](list-service-providers-output-service-discovery-information-properties-module-id.md))

## module\_uuid



`module_uuid`

* is required

* Type: `string` ([Module UUID](list-service-providers-output-service-discovery-information-properties-module-uuid.md))

* cannot be null

* defined in: [list-service-providers output](list-service-providers-output-service-discovery-information-properties-module-uuid.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/module_uuid")

### module\_uuid Type

`string` ([Module UUID](list-service-providers-output-service-discovery-information-properties-module-uuid.md))

## node

Node identifier where the module is hosted

`node`

* is required

* Type: `string` ([Node ID](list-service-providers-output-service-discovery-information-properties-node-id.md))

* cannot be null

* defined in: [list-service-providers output](list-service-providers-output-service-discovery-information-properties-node-id.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/node")

### node Type

`string` ([Node ID](list-service-providers-output-service-discovery-information-properties-node-id.md))

## transport



`transport`

* is required

* Type: `string` ([Transport protocol](list-service-providers-output-service-discovery-information-properties-transport-protocol.md))

* cannot be null

* defined in: [list-service-providers output](list-service-providers-output-service-discovery-information-properties-transport-protocol.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/transport")

### transport Type

`string` ([Transport protocol](list-service-providers-output-service-discovery-information-properties-transport-protocol.md))

### transport Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### transport Examples

```json
"tcp"
```

```json
"udp"
```

```json
"http"
```

## ui\_name

Custom UI label of the module

`ui_name`

* is required

* Type: `string` ([UI name](list-service-providers-output-service-discovery-information-properties-ui-name.md))

* can be null

* defined in: [list-service-providers output](list-service-providers-output-service-discovery-information-properties-ui-name.md "http://schema.nethserver.org/agent/list-service-providers-output.json#/itmes/properties/ui_name")

### ui\_name Type

`string` ([UI name](list-service-providers-output-service-discovery-information-properties-ui-name.md))
