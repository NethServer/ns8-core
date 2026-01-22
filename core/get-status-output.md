# Get module status Schema

```txt
http://schema.nethserver.org/agent/get-status-output.json
```

Get module instance running status

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-status-output.json](agent/get-status-output.json "open original schema") |

## Get module status Type

`object` ([Get module status](get-status-output.md))

## Get module status Examples

```json
{
  "instance": "dokuwiki1",
  "ui_name": "mywiki",
  "node": "1",
  "node_ui_name": "andromeda",
  "services": [
    {
      "name": "dokuwiki",
      "active": true,
      "enabled": true,
      "failed": false
    }
  ],
  "images": [
    {
      "name": "docker.io/bitnami/dokuwiki:20200729.0.0-debian-10-r299",
      "created": "2021-07-12 10:15:52 +0000 UTC",
      "size": "402 MB"
    }
  ],
  "volumes": [
    {
      "name": "dokuwiki-data",
      "mount": "/home/dokuwiki1/.local/share/containers/storage/volumes/dokuwiki-data/_data",
      "created": "2021-07-19 10:19:45.528366456 +0200 CEST"
    }
  ]
}
```

# Get module status Properties

| Property                        | Type     | Required | Nullable       | Defined by                                                                                                                                             |
| :------------------------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------- |
| [instance](#instance)           | `string` | Required | cannot be null | [Get module status](get-status-output-properties-instance.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/instance")         |
| [ui\_name](#ui_name)            | `string` | Required | cannot be null | [Get module status](get-status-output-properties-ui_name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/ui_name")           |
| [node](#node)                   | `string` | Required | cannot be null | [Get module status](get-status-output-properties-node.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/node")                 |
| [node\_ui\_name](#node_ui_name) | `string` | Required | cannot be null | [Get module status](get-status-output-properties-node_ui_name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/node_ui_name") |
| [services](#services)           | `array`  | Required | cannot be null | [Get module status](get-status-output-properties-services.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services")         |
| [images](#images)               | `array`  | Required | cannot be null | [Get module status](get-status-output-properties-images.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images")             |
| [volumes](#volumes)             | `array`  | Required | cannot be null | [Get module status](get-status-output-properties-volumes.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes")           |

## instance

Instance identifier

`instance`

* is required

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-instance.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/instance")

### instance Type

`string`

## ui\_name

UI name of the module instance

`ui_name`

* is required

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-ui_name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/ui_name")

### ui\_name Type

`string`

## node

Id of the node where the instance is running

`node`

* is required

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-node.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/node")

### node Type

`string`

## node\_ui\_name

UI name of the node

`node_ui_name`

* is required

* Type: `string`

* cannot be null

* defined in: [Get module status](get-status-output-properties-node_ui_name.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/node_ui_name")

### node\_ui\_name Type

`string`

## services

A list of systemd unit services

`services`

* is required

* Type: `object[]` ([Details](get-status-output-properties-services-items.md))

* cannot be null

* defined in: [Get module status](get-status-output-properties-services.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/services")

### services Type

`object[]` ([Details](get-status-output-properties-services-items.md))

## images

A list of podman images

`images`

* is required

* Type: `object[]` ([Details](get-status-output-properties-images-items.md))

* cannot be null

* defined in: [Get module status](get-status-output-properties-images.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/images")

### images Type

`object[]` ([Details](get-status-output-properties-images-items.md))

## volumes

A list of podman volumes

`volumes`

* is required

* Type: `object[]` ([Details](get-status-output-properties-volumes-items.md))

* cannot be null

* defined in: [Get module status](get-status-output-properties-volumes.md "http://schema.nethserver.org/agent/get-status-output.json#/properties/volumes")

### volumes Type

`object[]` ([Details](get-status-output-properties-volumes-items.md))
