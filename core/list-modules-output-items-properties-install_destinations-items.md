# Untitled object in list-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-modules-output.json\*](cluster/list-modules-output.json "open original schema") |

## items Type

`object` ([Details](list-modules-output-items-properties-install_destinations-items.md))

# items Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                               |
| :------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id)             | `integer` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-node_id.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/node_id")             |
| [instances](#instances)          | `integer` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-instances.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/instances")         |
| [eligible](#eligible)            | `boolean` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-eligible.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/eligible")           |
| [reject\_reason](#reject_reason) | `object`  | Required | can be null    | [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-reject_reason.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/reject_reason") |

## node\_id

Node identifier

`node_id`

* is required

* Type: `integer`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-node_id.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/node_id")

### node\_id Type

`integer`

## instances

Number of module instances currently installed on the node

`instances`

* is required

* Type: `integer`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-instances.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/instances")

### instances Type

`integer`

## eligible

True if another instance of the module can be installed on the node

`eligible`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-eligible.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/eligible")

### eligible Type

`boolean`

## reject\_reason



`reject_reason`

* is required

* Type: `object` ([Details](list-modules-output-items-properties-install_destinations-items-properties-reject_reason.md))

* can be null

* defined in: [list-modules output](list-modules-output-items-properties-install_destinations-items-properties-reject_reason.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations/items/properties/reject_reason")

### reject\_reason Type

`object` ([Details](list-modules-output-items-properties-install_destinations-items-properties-reject_reason.md))
