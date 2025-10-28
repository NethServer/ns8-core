# Untitled object in determine-restore-eligibility output Schema

```txt
http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [determine-restore-eligibility-output.json\*](cluster/determine-restore-eligibility-output.json "open original schema") |

## items Type

`object` ([Details](determine-restore-eligibility-output-properties-install_destinations-items.md))

# items Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                                      |
| :------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [node\_id](#node_id)             | `integer` | Required | cannot be null | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-node_id.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/node_id")             |
| [instances](#instances)          | `integer` | Required | cannot be null | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-instances.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/instances")         |
| [eligible](#eligible)            | `boolean` | Required | cannot be null | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-eligible.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/eligible")           |
| [reject\_reason](#reject_reason) | `object`  | Required | can be null    | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/reject_reason") |

## node\_id

Node identifier

`node_id`

* is required

* Type: `integer`

* cannot be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-node_id.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/node_id")

### node\_id Type

`integer`

## instances

Number of module instances currently installed on the node

`instances`

* is required

* Type: `integer`

* cannot be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-instances.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/instances")

### instances Type

`integer`

## eligible

True if another instance of the module can be installed on the node

`eligible`

* is required

* Type: `boolean`

* cannot be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-eligible.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/eligible")

### eligible Type

`boolean`

## reject\_reason



`reject_reason`

* is required

* Type: `object` ([Details](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason.md))

* can be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/reject_reason")

### reject\_reason Type

`object` ([Details](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason.md))
