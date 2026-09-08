# list-rich-rules output Schema

```txt
http://schema.nethserver.org/node/list-rich-rules-output.json
```

Output schema of the list-rich-rules action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-rich-rules-output.json](node/list-rich-rules-output.json "open original schema") |

## list-rich-rules output Type

`object` ([list-rich-rules output](list-rich-rules-output.md))

## list-rich-rules output Examples

```json
{
  "rich_rules": [
    "rule family=\"ipv4\" forward-port port=\"5060\" protocol=\"udp\" to-port=\"5060\"",
    "rule family=\"ipv4\" source address=\"10.1.2.3\" port port=\"22\" protocol=\"tcp\" accept"
  ]
}
```

# list-rich-rules output Properties

| Property                   | Type    | Required | Nullable       | Defined by                                                                                                                                                       |
| :------------------------- | :------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [rich\_rules](#rich_rules) | `array` | Required | cannot be null | [list-rich-rules output](list-rich-rules-output-properties-rich_rules.md "http://schema.nethserver.org/node/list-rich-rules-output.json#/properties/rich_rules") |

## rich\_rules



`rich_rules`

* is required

* Type: `string[]`

* cannot be null

* defined in: [list-rich-rules output](list-rich-rules-output-properties-rich_rules.md "http://schema.nethserver.org/node/list-rich-rules-output.json#/properties/rich_rules")

### rich\_rules Type

`string[]`
