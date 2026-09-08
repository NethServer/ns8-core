# Untitled undefined type in list-alerts output Schema

```txt
http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-alerts.json\*](cluster/list-alerts.json "open original schema") |

## 0 Type

unknown

# 0 Properties

| Property                 | Type          | Required | Nullable       | Defined by                                                                                                                                                                               |
| :----------------------- | :------------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [category](#category)    | Not specified | Optional | cannot be null | [list-alerts output](list-alerts-defs-base-attrs-oneof-0-properties-category.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/0/properties/category")   |
| [module\_id](#module_id) | `string`      | Required | cannot be null | [list-alerts output](list-alerts-defs-base-attrs-oneof-0-properties-module_id.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/0/properties/module_id") |

## category



`category`

* is optional

* Type: unknown

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-base-attrs-oneof-0-properties-category.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/0/properties/category")

### category Type

unknown

### category Constraints

**constant**: the value of this property must be equal to:

```json
"application"
```

## module\_id



`module_id`

* is required

* Type: `string`

* cannot be null

* defined in: [list-alerts output](list-alerts-defs-base-attrs-oneof-0-properties-module_id.md "http://schema.nethserver.org/cluster/list-alerts.json#/$defs/base-attrs/oneOf/0/properties/module_id")

### module\_id Type

`string`
