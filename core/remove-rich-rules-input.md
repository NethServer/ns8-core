# remove-rich-rules input Schema

```txt
http://schema.nethserver.org/node/remove-rich-rules-input.json
```

Remove firewall rich rules

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-rich-rules-input.json](node/remove-rich-rules-input.json "open original schema") |

## remove-rich-rules input Type

`object` ([remove-rich-rules input](remove-rich-rules-input.md))

## remove-rich-rules input Examples

```json
{
  "rich_rules": [
    "rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060",
    "rule family=ipv6 forward-port port=5060 protocol=udp to-port=5060 to-addr=2001:db8::1",
    "rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060 source address=1.2.3.4"
  ]
}
```

# remove-rich-rules input Properties

| Property                   | Type    | Required | Nullable       | Defined by                                                                                                                                                          |
| :------------------------- | :------ | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [rich\_rules](#rich_rules) | `array` | Required | cannot be null | [remove-rich-rules input](remove-rich-rules-input-properties-rich_rules.md "http://schema.nethserver.org/node/remove-rich-rules-input.json#/properties/rich_rules") |

## rich\_rules



`rich_rules`

* is required

* Type: `string[]`

* cannot be null

* defined in: [remove-rich-rules input](remove-rich-rules-input-properties-rich_rules.md "http://schema.nethserver.org/node/remove-rich-rules-input.json#/properties/rich_rules")

### rich\_rules Type

`string[]`

### rich\_rules Constraints

**minimum number of items**: the minimum number of items for this array is: `1`
