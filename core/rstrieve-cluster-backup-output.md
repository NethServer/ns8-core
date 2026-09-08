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
  "single_node_cluster": true,
  "cluster": "mycluster1",
  "vpn": "cs1.leader.cluster0.gs.nethserver.net:55820",
  "domains": 1,
  "backup_repositories": 1
}
```

# retrieve-cluster-backup output Properties

| Property                                      | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                              |
| :-------------------------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [single\_node\_cluster](#single_node_cluster) | `boolean` | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-single_node_cluster.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/single_node_cluster")           |
| [cluster](#cluster)                           | `string`  | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-cluster-label.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/cluster")                             |
| [vpn](#vpn)                                   | `string`  | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-vpn-endpoint.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/vpn")                                  |
| [domains](#domains)                           | `integer` | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-domains.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/domains")                         |
| [backup\_repositories](#backup_repositories)  | `integer` | Required | cannot be null | [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/backup_repositories") |

## single\_node\_cluster

If true, the backup refers to a cluster of a single node

`single_node_cluster`

* is required

* Type: `boolean`

* cannot be null

* defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-single_node_cluster.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/single_node_cluster")

### single\_node\_cluster Type

`boolean`

## cluster

User defined cluster label, can be empty'

`cluster`

* is required

* Type: `string` ([Cluster label](rstrieve-cluster-backup-output-properties-cluster-label.md))

* cannot be null

* defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-cluster-label.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/cluster")

### cluster Type

`string` ([Cluster label](rstrieve-cluster-backup-output-properties-cluster-label.md))

## vpn

Original VPN host and port

`vpn`

* is required

* Type: `string` ([VPN endpoint](rstrieve-cluster-backup-output-properties-vpn-endpoint.md))

* cannot be null

* defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-vpn-endpoint.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/vpn")

### vpn Type

`string` ([VPN endpoint](rstrieve-cluster-backup-output-properties-vpn-endpoint.md))

## domains

Number of user domains configured in the original cluster

`domains`

* is required

* Type: `integer` ([Number of domains](rstrieve-cluster-backup-output-properties-number-of-domains.md))

* cannot be null

* defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-domains.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/domains")

### domains Type

`integer` ([Number of domains](rstrieve-cluster-backup-output-properties-number-of-domains.md))

## backup\_repositories

Number of backup repositories configured in the original cluster

`backup_repositories`

* is required

* Type: `integer` ([Number of backup repositories](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md))

* cannot be null

* defined in: [retrieve-cluster-backup output](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md "http://schema.nethserver.org/cluster/rstrieve-cluster-backup-output.json#/properties/backup_repositories")

### backup\_repositories Type

`integer` ([Number of backup repositories](rstrieve-cluster-backup-output-properties-number-of-backup-repositories.md))
