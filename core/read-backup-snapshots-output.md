# read-backup-snapshots output Schema

```txt
http://schema.nethserver.org/cluster/read-backup-snapshots-output.json
```

Read the snaphost list of a given backup

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [read-backup-snapshots-output.json](cluster/read-backup-snapshots-output.json "open original schema") |

## read-backup-snapshots output Type

unknown ([read-backup-snapshots output](read-backup-snapshots-output.md))

## read-backup-snapshots output Examples

```json
[
  {
    "timestamp": 1643964838,
    "id": "a2a4cb238e4bd428900756376d4ff94009a8f487effe25c145bfaffa72406693"
  },
  {
    "timestamp": 1644037213,
    "id": "e993a67283eefa7a5d148284790e526f70f360788ed71e17e18406bc5a8a1185"
  },
  {
    "timestamp": 1644123618,
    "id": "4ea007b7770cfab125ab6035ae5371fe0e87464b8caa1607b2ced8c7e8732b4a"
  },
  {
    "timestamp": 1644210013,
    "id": "72a3cecc8b6acbd1610ebabd9aae1ac2b0b96864fa986a7437188036378e61ad"
  },
  {
    "timestamp": 1644296416,
    "id": "ebd99130107f282e47e7eb161de8fe14e0b832b7eacb04c9bf60761407e6081a"
  }
]
```
