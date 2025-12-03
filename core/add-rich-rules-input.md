# add-rich-rules input Schema

```txt
http://schema.nethserver.org/node/add-rich-rules-input.json
```

Add firewall rich rules

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                         |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-rich-rules-input.json](node/add-rich-rules-input.json "open original schema") |

## add-rich-rules input Type

`object` ([add-rich-rules input](add-rich-rules-input.md))

## add-rich-rules input Examples

```json
{
  "rich_rules": [
    "rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060",
    "rule family=ipv6 forward-port port=5060 protocol=udp to-port=5060 to-addr=2001:db8::1",
    "rule family=ipv4 forward-port port=5060 protocol=udp to-port=5060 source address=1.2.3.4"
  ]
}
```

# add-rich-rules input Properties

| Property                   | Type    | Required | Nullable       | Defined by                                                                                                                                                 |
| :------------------------- | :------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [rich\_rules](#rich_rules) | `array` | Required | cannot be null | [add-rich-rules input](add-rich-rules-input-properties-rich_rules.md "http://schema.nethserver.org/node/add-rich-rules-input.json#/properties/rich_rules") |

## rich\_rules



`rich_rules`

* is required

* Type: `string[]`

* cannot be null

* defined in: [add-rich-rules input](add-rich-rules-input-properties-rich_rules.md "http://schema.nethserver.org/node/add-rich-rules-input.json#/properties/rich_rules")

### rich\_rules Type

`string[]`

### rich\_rules Constraints

**minimum number of items**: the minimum number of items for this array is: `1`
