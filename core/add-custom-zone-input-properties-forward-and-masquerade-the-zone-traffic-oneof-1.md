# Untitled object in add-custom-zone input Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [add-custom-zone-input.json\*](node/add-custom-zone-input.json "open original schema") |

## 1 Type

`object` ([Details](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1.md))

# 1 Properties

| Property                      | Type      | Required | Nullable       | Defined by                                                                                                                                                                                                                                                         |
| :---------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [internet](#internet)         | `boolean` | Optional | cannot be null | [add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1-properties-allow-the-internet.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1/properties/internet")       |
| [destinations](#destinations) | `array`   | Optional | cannot be null | [add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1-properties-allowed-destinations.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1/properties/destinations") |

## internet

Private networks stay closed unless a destination lists them

`internet`

* is optional

* Type: `boolean` ([Allow the Internet](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1-properties-allow-the-internet.md))

* cannot be null

* defined in: [add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1-properties-allow-the-internet.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1/properties/internet")

### internet Type

`boolean` ([Allow the Internet](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1-properties-allow-the-internet.md))

### internet Default Value

The default value is:

```json
false
```

## destinations

IPv4 address or network, with an optional port and protocol

`destinations`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1-properties-allowed-destinations.md "http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade/oneOf/1/properties/destinations")

### destinations Type

`string[]`

### destinations Examples

```json
"10.5.5.0/24:443/tcp"
```

```json
"192.168.1.10:445/tcp"
```

```json
"192.168.1.0/24"
```
