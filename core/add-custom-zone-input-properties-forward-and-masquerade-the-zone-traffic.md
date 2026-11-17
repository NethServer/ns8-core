# Forward and masquerade the zone traffic Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/masquerade
```

Forward the traffic of the zone with the node address as source. The interface must be a tun created by add-tun. The policy name is the zone name with the -masq suffix, so the zone name is limited to 13 characters. true is the same as {"internet": true}.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-custom-zone-input.json\*](node/add-custom-zone-input.json "open original schema") |

## masquerade Type

merged type ([Forward and masquerade the zone traffic](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic.md))

one (and only one) of

* [Untitled boolean in add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-0.md "check type definition")

* [Untitled object in add-custom-zone input](add-custom-zone-input-properties-forward-and-masquerade-the-zone-traffic-oneof-1.md "check type definition")

## masquerade Default Value

The default value is:

```json
false
```
