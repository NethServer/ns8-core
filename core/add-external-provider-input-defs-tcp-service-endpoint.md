# TCP service endpoint Schema

```txt
http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint
```

Initial TCP backend endpoint configuration

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-external-provider-input.json\*](cluster/add-external-provider-input.json "open original schema") |

## tcp-service-endpoint Type

`object` ([TCP service endpoint](add-external-provider-input-defs-tcp-service-endpoint.md))

# tcp-service-endpoint Properties

| Property      | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                  |
| :------------ | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [host](#host) | `string`  | Required | cannot be null | [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/host") |
| [port](#port) | `integer` | Required | cannot be null | [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/port") |

## host



`host`

* is required

* Type: `string`

* cannot be null

* defined in: [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-host.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/host")

### host Type

`string`

## port



`port`

* is required

* Type: `integer`

* cannot be null

* defined in: [add-external-provider input](add-external-provider-input-defs-tcp-service-endpoint-properties-port.md "http://schema.nethserver.org/cluster/add-external-provider-input.json#/$defs/tcp-service-endpoint/properties/port")

### port Type

`integer`

### port Constraints

**maximum**: the value of this number must smaller than or equal to: `65535`

**minimum**: the value of this number must greater than or equal to: `1`
