# add-custom-zone input Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json
```

Add firewall configuration for a custom zone

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                           |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-custom-zone-input.json](node/add-custom-zone-input.json "open original schema") |

## add-custom-zone input Type

`object` ([add-custom-zone input](add-custom-zone-input.md))

## add-custom-zone input Examples

```json
{
  "zone": "myzone",
  "interface": "eth0",
  "ports": [
    "25/tcp"
  ],
  "rules": [
    "rule family=ipv4 source address=172.18.212.0/24 destination address=172.18.212.0/24 reject"
  ]
}
```

```json
{
  "zone": "tunvpn1",
  "interface": "tunvpn1",
  "ports": [
    "53/udp"
  ],
  "rules": [
    "rule family=ipv4 source address=10.8.0.0/24 port port=22 protocol=tcp reject"
  ],
  "masquerade": {
    "internet": true,
    "destinations": [
      "10.5.5.0/24:443/tcp"
    ]
  }
}
```

# add-custom-zone input Properties

| Property                  | Type     | Required | Nullable       | Defined by                                                                                                                                                                                   |
| :------------------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [zone](#zone)             | `string` | Required | cannot be null | [add-custom-zone input](add-custom-zone-input-properties-zone-name.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/zone")                                       |
| [interface](#interface)   | `string` | Required | cannot be null | [add-custom-zone input](add-custom-zone-input-properties-interface-name-to-be-added-inside-the-zone.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/interface") |
| [ports](#ports)           | `array`  | Required | cannot be null | [add-custom-zone input](add-custom-zone-input-properties-firewall-ports-configuration.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/ports")                   |
| [rules](#rules)           | `array`  | Optional | cannot be null | [add-custom-zone input](add-custom-zone-input-properties-firewall-rich-rules.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/rules")                            |
| [masquerade](#masquerade) | Merged   | Optional | cannot be null | [add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade")   |

## zone



`zone`

* is required

* Type: `string` ([Zone name](add-custom-zone-input-properties-zone-name.md))

* cannot be null

* defined in: [add-custom-zone input](add-custom-zone-input-properties-zone-name.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/zone")

### zone Type

`string` ([Zone name](add-custom-zone-input-properties-zone-name.md))

### zone Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## interface



`interface`

* is required

* Type: `string` ([Interface name to be added inside the zone](add-custom-zone-input-properties-interface-name-to-be-added-inside-the-zone.md))

* cannot be null

* defined in: [add-custom-zone input](add-custom-zone-input-properties-interface-name-to-be-added-inside-the-zone.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/interface")

### interface Type

`string` ([Interface name to be added inside the zone](add-custom-zone-input-properties-interface-name-to-be-added-inside-the-zone.md))

### interface Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## ports



`ports`

* is required

* Type: `array` ([Firewall ports configuration](add-custom-zone-input-properties-firewall-ports-configuration.md))

* cannot be null

* defined in: [add-custom-zone input](add-custom-zone-input-properties-firewall-ports-configuration.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/ports")

### ports Type

`array` ([Firewall ports configuration](add-custom-zone-input-properties-firewall-ports-configuration.md))

### ports Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

### ports Examples

```json
"25/tcp"
```

```json
"21/udp"
```

```json
"9000-9010/tcp"
```

## rules



`rules`

* is optional

* Type: `array` ([Firewall rich rules](add-custom-zone-input-properties-firewall-rich-rules.md))

* cannot be null

* defined in: [add-custom-zone input](add-custom-zone-input-properties-firewall-rich-rules.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/rules")

### rules Type

`array` ([Firewall rich rules](add-custom-zone-input-properties-firewall-rich-rules.md))

### rules Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

### rules Examples

```json
"rule family=ipv4 source address=172.18.212.1 destination address=172.18.212.0/24 accept"
```

## masquerade

Forward the traffic of the zone with the node address as source. The interface must be a tun created by add-tun. The policy name is the zone name with the -masq suffix, so the zone name is limited to 13 characters. true is the same as {"internet": true}.

`masquerade`

* is optional

* Type: merged type ([Forward and masquerade the zone traffic](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic.md))

* cannot be null

* defined in: [add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade")

### masquerade Type

merged type ([Forward and masquerade the zone traffic](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic.md))

one (and only one) of

* [Untitled boolean in add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-0.md "check type definition")

* [Untitled object in add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1.md "check type definition")

### masquerade Default Value

The default value is:

```json
false
```
