# add-rclone-destination input Schema

```txt
http://schema.nethserver.org/node/add-rclone-destination-input.json
```

Input schema of the add-rclone-destination action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                         |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-rclone-destination-input.json](node/add-rclone-destination-input.json "open original schema") |

## add-rclone-destination input Type

`object` ([add-rclone-destination input](add-rclone-destination-input.md))

## add-rclone-destination input Examples

```json
{
  "id": "8d3a9df9-f20a-4d4c-9280-718126191d8f",
  "basepath": "mybucket",
  "rclone_conf_secret": "[8d3a9df9-f20a-4d4c-9280-718126191d8f]\n..."
}
```

# add-rclone-destination input Properties

| Property                                    | Type     | Required | Nullable       | Defined by                                                                                                                                                                                         |
| :------------------------------------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)                                   | `string` | Required | cannot be null | [add-rclone-destination input](add-rclone-destination-input-properties-id.md "http://schema.nethserver.org/node/add-rclone-destination-input.json#/properties/id")                                 |
| [basepath](#basepath)                       | `string` | Required | cannot be null | [add-rclone-destination input](add-rclone-destination-input-properties-basepath.md "http://schema.nethserver.org/node/add-rclone-destination-input.json#/properties/basepath")                     |
| [rclone\_conf\_secret](#rclone_conf_secret) | `string` | Required | cannot be null | [add-rclone-destination input](add-rclone-destination-input-properties-rclone_conf_secret.md "http://schema.nethserver.org/node/add-rclone-destination-input.json#/properties/rclone_conf_secret") |

## id



`id`

* is required

* Type: `string`

* cannot be null

* defined in: [add-rclone-destination input](add-rclone-destination-input-properties-id.md "http://schema.nethserver.org/node/add-rclone-destination-input.json#/properties/id")

### id Type

`string`

## basepath



`basepath`

* is required

* Type: `string`

* cannot be null

* defined in: [add-rclone-destination input](add-rclone-destination-input-properties-basepath.md "http://schema.nethserver.org/node/add-rclone-destination-input.json#/properties/basepath")

### basepath Type

`string`

## rclone\_conf\_secret



`rclone_conf_secret`

* is required

* Type: `string`

* cannot be null

* defined in: [add-rclone-destination input](add-rclone-destination-input-properties-rclone_conf_secret.md "http://schema.nethserver.org/node/add-rclone-destination-input.json#/properties/rclone_conf_secret")

### rclone\_conf\_secret Type

`string`
