# get-node-status output Schema

```txt
http://schema.nethserver.org/node/get-node-status-output.json
```

Output schema of the get-node-status action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-node-status-output.json](node/get-node-status-output.json "open original schema") |

## get-node-status output Type

`object` ([get-node-status output](get-node-status-output.md))

## get-node-status output Examples

```json
{
  "load": {
    "1min": 0,
    "5min": 0,
    "15min": 0
  },
  "cpu": {
    "count": 1,
    "usage": 0,
    "info": [
      {
        "vendor": "GenuineIntel",
        "model": "DO-Premium-Intel",
        "clock": "2494.144"
      }
    ]
  },
  "memory": {
    "total": 1023934464,
    "used": 370835456,
    "free": 653099008
  },
  "swap": {
    "total": 0,
    "used": 0,
    "free": 0
  },
  "disks": [
    {
      "device": "/dev/vda1",
      "mountpoint": "/",
      "fstype": "ext4",
      "total": 26638249984,
      "used": 2510962688,
      "free": 23022039040
    },
    {
      "device": "/dev/vda15",
      "mountpoint": "/boot/efi",
      "fstype": "vfat",
      "total": 129718272,
      "used": 6160384,
      "free": 123557888
    },
    {
      "device": "/dev/vda1",
      "mountpoint": "/var/lib/containers/storage/overlay",
      "fstype": "ext4",
      "total": 26638249984,
      "used": 2510962688,
      "free": 23022039040
    }
  ]
}
```

# get-node-status output Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                               |
| :---------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [load](#load)     | `object` | Required | cannot be null | [get-node-status output](get-node-status-output-properties-load.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/load")     |
| [cpu](#cpu)       | `object` | Required | cannot be null | [get-node-status output](get-node-status-output-properties-cpu.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu")       |
| [memory](#memory) | `object` | Optional | cannot be null | [get-node-status output](get-node-status-output-properties-memory.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory") |
| [swap](#swap)     | `object` | Required | cannot be null | [get-node-status output](get-node-status-output-properties-swap.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap")     |
| [disks](#disks)   | `array`  | Required | cannot be null | [get-node-status output](get-node-status-output-properties-disks.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks")   |

## load

Average load

`load`

* is required

* Type: `object` ([Details](get-node-status-output-properties-load.md))

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-load.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/load")

### load Type

`object` ([Details](get-node-status-output-properties-load.md))

## cpu

Average load

`cpu`

* is required

* Type: `object` ([Details](get-node-status-output-properties-cpu.md))

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-cpu.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/cpu")

### cpu Type

`object` ([Details](get-node-status-output-properties-cpu.md))

## memory

Memory statistics

`memory`

* is optional

* Type: `object` ([Details](get-node-status-output-properties-memory.md))

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-memory.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/memory")

### memory Type

`object` ([Details](get-node-status-output-properties-memory.md))

## swap

SWAP statistics

`swap`

* is required

* Type: `object` ([Details](get-node-status-output-properties-swap.md))

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-swap.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/swap")

### swap Type

`object` ([Details](get-node-status-output-properties-swap.md))

## disks

Storage information by partition

`disks`

* is required

* Type: `object[]` ([Details](get-node-status-output-properties-disks-items.md))

* cannot be null

* defined in: [get-node-status output](get-node-status-output-properties-disks.md "http://schema.nethserver.org/node/get-node-status-output.json#/properties/disks")

### disks Type

`object[]` ([Details](get-node-status-output-properties-disks-items.md))
