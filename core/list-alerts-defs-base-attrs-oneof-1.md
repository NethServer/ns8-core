# Untitled undefined type in list-alerts output Schema

```txt
http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/1
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-alerts.json\*](cluster/list-alerts.json "open original schema") |

## 1 Type

unknown

# 1 Properties

| Property              | Type          | Required | Nullable       | Defined by                                                                                                                                                                             |
| :-------------------- | :------------ | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [category](#category) | Not specified | Optional | cannot be null | [list-alerts output](list-alerts-defs-base-attrs-oneof-1-properties-category.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/1/properties/category") |
| [node\_id](#node_id)  | `integer`     | Required | cannot be null | [list-alerts output](list-alerts-defs-base-attrs-oneof-1-properties-node_id.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/1/properties/node_id")   |

## category



`category`

* is optional

* Type: unknown

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-base-attrs-oneof-1-properties-category.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/1/properties/category")

### category Type

unknown

### category Constraints

**constant**: the value of this property must be equal to:

```json
"node"
```

## node\_id



`node_id`

* is required

* Type: `integer`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-base-attrs-oneof-1-properties-node_id.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/1/properties/node_id")

### node\_id Type

`integer`
