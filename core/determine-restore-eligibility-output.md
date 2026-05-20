# determine-restore-eligibility output Schema

```txt
http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json
```

Output schema of the determine-restore-eligibility action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [determine-restore-eligibility-output.json](cluster/determine-restore-eligibility-output.json "open original schema") |

## determine-restore-eligibility output Type

`object` ([determine-restore-eligibility output](determine-restore-eligibility-output.md))

## determine-restore-eligibility output Examples

```json
{
  "image_url": "ghcr.io/nethserver/ejabberd:1.0.1",
  "install_destinations": [
    {
      "node_id": 1,
      "instances": 0,
      "eligible": true,
      "reject_reason": null
    },
    {
      "node_id": 2,
      "instances": 1,
      "eligible": false,
      "reject_reason": {
        "message": "max_per_node_limit",
        "parameter": "1"
      }
    }
  ]
}
```

# determine-restore-eligibility output Properties

| Property                                       | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                                        |
| :--------------------------------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [image\_url](#image_url)                       | `string` | Required | cannot be null | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-image_url.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/image_url")                       |
| [install\_destinations](#install_destinations) | `array`  | Required | cannot be null | [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations") |

## image\_url



`image_url`

* is required

* Type: `string`

* cannot be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-image_url.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/image_url")

### image\_url Type

`string`

## install\_destinations

Describe for each node of the cluster if the node is eligible or not to install a new module instance. If not, a reject reason is returned.

`install_destinations`

* is required

* Type: `object[]` ([Details](determine-restore-eligibility-output-properties-install_destinations-items.md))

* cannot be null

* defined in: [determine-restore-eligibility output](determine-restore-eligibility-output-properties-install_destinations.md "http://schema.nethserver.org/cluster/determine-restore-eligibility-output.json#/properties/install_destinations")

### install\_destinations Type

`object[]` ([Details](determine-restore-eligibility-output-properties-install_destinations-items.md))
