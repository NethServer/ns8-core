# get metrics id Schema

```txt
http://schema.nethserver.org/cluster/get-smarthost.json
```

Get the module id of the default Metrics instance

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-smarthost.json](cluster/get-smarthost.json "open original schema") |

## get metrics id Type

`object` ([get metrics id](get-smarthost.md))

## get metrics id Examples

```json
{
  "metrics_id": "metrics1"
}
```

# get metrics id Properties

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                |
| :------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------- |
| [metrics\_id](#metrics_id) | `string` | Required | cannot be null | [get metrics id](get-smarthost-properties-metrics_id.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/metrics_id") |

## metrics\_id

The ID of the default metrics instance

`metrics_id`

* is required

* Type: `string` ([metrics\_id](get-smarthost-properties-metrics_id.md))

* cannot be null

* defined in: [get metrics id](get-smarthost-properties-metrics_id.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/metrics_id")

### metrics\_id Type

`string` ([metrics\_id](get-smarthost-properties-metrics_id.md))
