# retrieve-cluster-backup output Schema

```txt
http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json
```

Return info from cluster backup

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [rstrieve-cluster-backup-output.json](cluster/rstrieve-cluster-backup-output.json "open original schema") |

## retrieve-cluster-backup output Type

`object` ([retrieve-cluster-backup output](rstrieve-cluster-backup-output.md))

## retrieve-cluster-backup output Examples

```json
{
  "timestamp": 1644230044,
  "cluster": "mycluster1",
  "vpn": "cs1.leader.cluster0.gs.nethserver.net:55820",
  "domains": 1,
  "backup_repositories": 1
}
```

# retrieve-cluster-backup output Properties

| Property                                     | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                              |
| :------------------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [timestamp](#timestamp)                      | `integer` | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-backup-timestamp.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/timestamp")                        |
| [cluster](#cluster)                          | `string`  | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-cluster-label.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/cluster")                             |
| [vpn](#vpn)                                  | `string`  | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-vpn-endpoint.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/vpn")                                  |
| [domains](#domains)                          | `integer` | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-domains.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/domains")                         |
| [backup\_repositories](#backup_repositories) | `integer` | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/backup_repositories") |

## timestamp

Date and time of the backup, timestamp in UTC timezone

`timestamp`

*   is required

*   Type: `integer` ([Backup timestamp](rstrieve-cluster-backup-output-properties-backup-timestamp.md))

*   cannot be null

*   defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-backup-timestamp.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/timestamp")

### timestamp Type

`integer` ([Backup timestamp](rstrieve-cluster-backup-output-properties-backup-timestamp.md))

## cluster

User defined cluster label, can be empty'

`cluster`

*   is required

*   Type: `string` ([Cluster label](rstrieve-cluster-backup-output-properties-cluster-label.md))

*   cannot be null

*   defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-cluster-label.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/cluster")

### cluster Type

`string` ([Cluster label](rstrieve-cluster-backup-output-properties-cluster-label.md))

## vpn

Original VPN host and port

`vpn`

*   is required

*   Type: `string` ([VPN endpoint](rstrieve-cluster-backup-output-properties-vpn-endpoint.md))

*   cannot be null

*   defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-vpn-endpoint.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/vpn")

### vpn Type

`string` ([VPN endpoint](rstrieve-cluster-backup-output-properties-vpn-endpoint.md))

## domains

Number of user domains configured in the original cluster

`domains`

*   is required

*   Type: `integer` ([Number of domains](rstrieve-cluster-backup-output-properties-number-of-domains.md))

*   cannot be null

*   defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-domains.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/domains")

### domains Type

`integer` ([Number of domains](rstrieve-cluster-backup-output-properties-number-of-domains.md))

## backup\_repositories

Number of backup repositories configured in the original cluster

`backup_repositories`

*   is required

*   Type: `integer` ([Number of backup repositories](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md))

*   cannot be null

*   defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/backup_repositories")

### backup\_repositories Type

`integer` ([Number of backup repositories](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md))
