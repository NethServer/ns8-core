# transfer-state input Schema

```txt
http://schema.nethserver.org/module/transfer-state-input.json
```

Transfer the module state to another module instance

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                           |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [transfer-state-input.json](module/transfer-state-input.json "open original schema") |

## transfer-state input Type

`object` ([transfer-state input](transfer-state-input.md))

## transfer-state input Examples

```json
{
  "replace": true,
  "credentials": [
    "dokuwiki1",
    "s3cr3t"
  ],
  "address": "10.5.4.3",
  "port": 20021
}
```

# transfer-state input Properties

| Property                    | Type      | Required | Nullable       | Defined by                                                                                                                                                                    |
| :-------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [replace](#replace)         | `boolean` | Required | cannot be null | [transfer-state input](transfer-state-input-properties-replace-flag.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/replace")                   |
| [credentials](#credentials) | `array`   | Required | cannot be null | [transfer-state input](transfer-state-input-properties-rsyncd-service-credentials.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/credentials") |
| [address](#address)         | `string`  | Required | cannot be null | [transfer-state input](transfer-state-input-properties-rsyncd-remote-address.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/address")          |
| [port](#port)               | `integer` | Required | cannot be null | [transfer-state input](transfer-state-input-properties-rsyncd-remote-tcp-port-number.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/port")     |

## replace

If set to true the state cannot be modified after the transfer is completed

`replace`

* is required

* Type: `boolean` ([Replace flag](transfer-state-input-properties-replace-flag.md))

* cannot be null

* defined in: [transfer-state input](transfer-state-input-properties-replace-flag.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/replace")

### replace Type

`boolean` ([Replace flag](transfer-state-input-properties-replace-flag.md))

## credentials

An array with two elements: username, password

`credentials`

* is required

* Type: `string[]`

* cannot be null

* defined in: [transfer-state input](transfer-state-input-properties-rsyncd-service-credentials.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/credentials")

### credentials Type

`string[]`

### credentials Constraints

**maximum number of items**: the maximum number of items for this array is: `2`

**minimum number of items**: the minimum number of items for this array is: `2`

## address

It must be an host name or IP address

`address`

* is required

* Type: `string` ([Rsyncd remote address](transfer-state-input-properties-rsyncd-remote-address.md))

* cannot be null

* defined in: [transfer-state input](transfer-state-input-properties-rsyncd-remote-address.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/address")

### address Type

`string` ([Rsyncd remote address](transfer-state-input-properties-rsyncd-remote-address.md))

### address Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## port



`port`

* is required

* Type: `integer` ([Rsyncd remote TCP port number](transfer-state-input-properties-rsyncd-remote-tcp-port-number.md))

* cannot be null

* defined in: [transfer-state input](transfer-state-input-properties-rsyncd-remote-tcp-port-number.md "http://schema.nethserver.org/module/transfer-state-input.json#/properties/port")

### port Type

`integer` ([Rsyncd remote TCP port number](transfer-state-input-properties-rsyncd-remote-tcp-port-number.md))

### port Constraints

**minimum**: the value of this number must greater than or equal to: `1`
