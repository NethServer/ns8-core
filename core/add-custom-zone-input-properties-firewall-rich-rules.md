# Firewall rich rules Schema

```txt
http://schema.nethserver.org/node/add-custom-zone-input.json#/properties/rules
```



| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                             |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-custom-zone-input.json\*](node/add-custom-zone-input.json "open original schema") |

## rules Type

`array` ([Firewall rich rules](add-custom-zone-input-properties-firewall-rich-rules.md))

## rules Constraints

**minimum number of items**: the minimum number of items for this array is: `1`

## rules Examples

```json
"rule family=ipv4 source address=172.18.212.1 destination address=172.18.212.0/24 accept"
```
