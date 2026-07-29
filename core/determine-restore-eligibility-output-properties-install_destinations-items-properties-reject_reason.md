# Untitled undefined type in determine-restore-eligibility output Schema

```txt
http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/reject_reason
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [determine-restore-eligibility-output.json\*](cluster/determine-restore-eligibility-output.json "open original schema") |

## reject\_reason Type

`object` ([Details](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason.md))

# reject\_reason Properties

| Property                | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                                                                                                                                |
| :---------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [message](#message)     | `string` | Required | cannot be null | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason-properties-message.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/reject_reason/properties/message")     |
| [parameter](#parameter) | `string` | Required | cannot be null | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason-properties-parameter.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/reject_reason/properties/parameter") |

## message



`message`

* is required

* Type: `string`

* cannot be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason-properties-message.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/reject_reason/properties/message")

### message Type

`string`

## parameter



`parameter`

* is required

* Type: `string`

* cannot be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations-items-properties-reject_reason-properties-parameter.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations/items/properties/reject_reason/properties/parameter")

### parameter Type

`string`
