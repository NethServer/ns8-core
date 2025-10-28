# remove-port-forward input Schema

```txt
http://schema.nethserver.org/node/remove-port-forward-input.json
```

Remove firewall configuration for a port forward

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                   |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-port-forward-input.json](node/remove-port-forward-input.json "open original schema") |

## remove-port-forward input Type

`object` ([remove-port-forward input](remove-port-forward-input.md))

## remove-port-forward input Examples

```json
{
  "external_port": "5060",
  "internal_port": "5060",
  "protocols": [
    "udp",
    "tcp"
  ]
}
```

```json
{
  "external_port": "4422",
  "internal_port": "22",
  "protocols": [
    "udp",
    "tcp"
  ],
  "destination_ip": "192.168.1.100"
}
```

# remove-port-forward input Properties

| Property                           | Type     | Required | Nullable       | Defined by                                                                                                                                                                                    |
| :--------------------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [external\_port](#external_port)   | `string` | Required | cannot be null | [remove-port-forward input](remove-port-forward-input-properties-external-port-or-port-range.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/external_port") |
| [internal\_port](#internal_port)   | `string` | Required | cannot be null | [remove-port-forward input](remove-port-forward-input-properties-internal-port-or-port-range.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/internal_port") |
| [protocols](#protocols)            | `array`  | Required | cannot be null | [remove-port-forward input](remove-port-forward-input-properties-protocols.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/protocols")                       |
| [destination\_ip](#destination_ip) | `string` | Optional | cannot be null | [remove-port-forward input](remove-port-forward-input-properties-destination-ip-address.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/destination_ip")     |

## external\_port



`external_port`

* is required

* Type: `string` ([External port or port range](remove-port-forward-input-properties-external-port-or-port-range.md))

* cannot be null

* defined in: [remove-port-forward input](remove-port-forward-input-properties-external-port-or-port-range.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/external_port")

### external\_port Type

`string` ([External port or port range](remove-port-forward-input-properties-external-port-or-port-range.md))

### external\_port Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## internal\_port



`internal_port`

* is required

* Type: `string` ([Internal port or port range](remove-port-forward-input-properties-internal-port-or-port-range.md))

* cannot be null

* defined in: [remove-port-forward input](remove-port-forward-input-properties-internal-port-or-port-range.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/internal_port")

### internal\_port Type

`string` ([Internal port or port range](remove-port-forward-input-properties-internal-port-or-port-range.md))

### internal\_port Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## protocols



`protocols`

* is required

* Type: `string[]`

* cannot be null

* defined in: [remove-port-forward input](remove-port-forward-input-properties-protocols.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/protocols")

### protocols Type

`string[]`

### protocols Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

## destination\_ip



`destination_ip`

* is optional

* Type: `string` ([Destination IP address](remove-port-forward-input-properties-destination-ip-address.md))

* cannot be null

* defined in: [remove-port-forward input](remove-port-forward-input-properties-destination-ip-address.md "http://schema.nethserver.org/node/remove-port-forward-input.json#/properties/destination_ip")

### destination\_ip Type

`string` ([Destination IP address](remove-port-forward-input-properties-destination-ip-address.md))

### destination\_ip Constraints

**IPv4**: the string must be an IPv4 address (dotted quad), according to [RFC 2673, section 3.2](https://tools.ietf.org/html/rfc2673 "check the specification")
