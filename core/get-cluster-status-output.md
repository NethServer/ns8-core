# get-cluster-status output Schema

```txt
http://schema.nethserver.org/cluster/get-cluster-status-output.json
```

Output schema of the get-cluster-status action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-cluster-status-output.json](cluster/get-cluster-status-output.json "open original schema") |

## get-cluster-status output Type

`object` ([get-cluster-status output](get-cluster-status-output.md))

## get-cluster-status output Examples

```json
{
  "ui_name": "mycluster",
  "initialized": true,
  "leader": true,
  "nodes": [
    {
      "id": 2,
      "ui_name": "node2@ams3",
      "local": false,
      "online": true,
      "hostname": "dn2.worker.cluster0.gs.nethserver.net",
      "machine_id": "96b91e27d15ae6d7ef0cc79c621f7904",
      "vpn": {
        "ip_address": "10.5.4.2",
        "public_key": "0yLfAut1CwatFb0k0sim6oPudIOWMzJsO9ggYVmvXBg=",
        "endpoint": "1.2.3.4",
        "last_seen": 1631198401,
        "rcvd": 5989132,
        "sent": 6097216,
        "keepalive": 25
      }
    },
    {
      "id": 1,
      "ui_name": "node1@sfo3",
      "local": true,
      "online": false,
      "vpn": {
        "ip_address": "10.5.4.1",
        "public_key": "c4BSbQV3zG1qIK1eSUfVvwLtTBxRuI1y5IWmddZgals="
      }
    }
  ],
  "leader_url": "https://server1.nethserver.org/cluster-admin/",
  "default_password": true
}
```

# get-cluster-status output Properties

| Property                               | Type      | Required | Nullable       | Defined by                                                                                                                                                                               |
| :------------------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [initialized](#initialized)            | `boolean` | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-initialized.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/initialized")           |
| [ui\_name](#ui_name)                   | `string`  | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-display-name.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/ui_name")              |
| [leader](#leader)                      | `boolean` | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-leader.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/leader")                     |
| [leader\_url](#leader_url)             | `string`  | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-leader_url.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/leader_url")             |
| [default\_password](#default_password) | `boolean` | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-default_password.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/default_password") |
| [nodes](#nodes)                        | `array`   | Required | cannot be null | [get-cluster-status output](get-cluster-status-output-properties-nodes.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes")                       |

## initialized

Set to true if the node has been configured, false otherwise

`initialized`

* is required

* Type: `boolean`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-initialized.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/initialized")

### initialized Type

`boolean`

## ui\_name



`ui_name`

* is required

* Type: `string` ([Display name](get-cluster-status-output-properties-display-name.md))

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-display-name.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/ui_name")

### ui\_name Type

`string` ([Display name](get-cluster-status-output-properties-display-name.md))

## leader

Set to true if the running node is the leader, false otherwise

`leader`

* is required

* Type: `boolean`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-leader.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/leader")

### leader Type

`boolean`

## leader\_url

UI URL of leader node

`leader_url`

* is required

* Type: `string`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-leader_url.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/leader_url")

### leader\_url Type

`string`

## default\_password

Set to true if the admin user has still the default passwsord, false otherwise

`default_password`

* is required

* Type: `boolean`

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-default_password.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/default_password")

### default\_password Type

`boolean`

## nodes



`nodes`

* is required

* Type: `object[]` ([Details](get-cluster-status-output-properties-nodes-items.md))

* cannot be null

* defined in: [get-cluster-status output](get-cluster-status-output-properties-nodes.md "http://schema.nethserver.org/cluster/get-cluster-status-output.json#/properties/nodes")

### nodes Type

`object[]` ([Details](get-cluster-status-output-properties-nodes-items.md))
