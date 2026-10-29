# add-tun input Schema

```txt
http://schema.nethserver.org/node/add-tun-input.json
```

Add a tun interface

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                           |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-tun-input.json](node/add-tun-input.json "open original schema") |

## add-tun input Type

`object` ([add-tun input](add-tun-input.md))

## add-tun input Examples

```json
{
  "tun": "tun0",
  "ip": "192.168.100.1/255.255.255.0"
}
```

```json
{
  "tun": "tun1",
  "ip": "192.168.200.1/24"
}
```

# add-tun input Properties

| Property    | Type     | Required | Nullable       | Defined by                                                                                                                                 |
| :---------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------- |
| [tun](#tun) | `string` | Required | cannot be null | [add-tun input](add-tun-input-properties-tun-name.md "http://schema.nethserver.org/node/add-tun-input.json#/properties/tun")               |
| [ip](#ip)   | `string` | Required | cannot be null | [add-tun input](add-tun-input-properties-ip-address-with-network.md "http://schema.nethserver.org/node/add-tun-input.json#/properties/ip") |

## tun



`tun`

* is required

* Type: `string` ([TUN name](add-tun-input-properties-tun-name.md))

* cannot be null

* defined in: [add-tun input](add-tun-input-properties-tun-name.md "http://schema.nethserver.org/node/add-tun-input.json#/properties/tun")

### tun Type

`string` ([TUN name](add-tun-input-properties-tun-name.md))

### tun Constraints

**minimum length**: the minimum number of characters for this string is: `4`

## ip



`ip`

* is required

* Type: `string` ([IP address with network](add-tun-input-properties-ip-address-with-network.md))

* cannot be null

* defined in: [add-tun input](add-tun-input-properties-ip-address-with-network.md "http://schema.nethserver.org/node/add-tun-input.json#/properties/ip")

### ip Type

`string` ([IP address with network](add-tun-input-properties-ip-address-with-network.md))

### ip Constraints

**minimum length**: the minimum number of characters for this string is: `1`
