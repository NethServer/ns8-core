# Sorting rule Schema

```txt
http://schema.nethserver.org/agent/list-service-providers-input.json#/properties/sortby
```

Sort the returned entries with the given rule

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-service-providers-input.json\*](agent/list-service-providers-input.json "open original schema") |

## sortby Type

`string` ([Sorting rule](list-service-providers-input-properties-sorting-rule.md))

## sortby Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value               | Explanation |
| :------------------ | :---------- |
| `"node_proximity"`  |             |
| `"module_seq_asc"`  |             |
| `"module_seq_desc"` |             |

## sortby Default Value

The default value is:

```json
"node_proximity"
```
