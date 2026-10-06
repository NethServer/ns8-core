# list-alerts input Schema

```txt
http://schema.nethserver.org/cluster/list-alerts-input.json
```

Input schema of the list-alerts action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-alerts-input.json](cluster/list-alerts-input.json "open original schema") |

## list-alerts input Type

`object` ([list-alerts input](list-alerts-input.md))

# list-alerts input Properties

| Property                  | Type    | Required | Nullable       | Defined by                                                                                                                                           |
| :------------------------ | :------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
| [categories](#categories) | `array` | Optional | cannot be null | [list-alerts input](list-alerts-input-properties-categories.md "http://schema.nethserver.org/cluster/list-alerts-input.json#/properties/categories") |

## categories



`categories`

* is optional

* Type: `string[]`

* cannot be null

* defined in: [list-alerts input](list-alerts-input-properties-categories.md "http://schema.nethserver.org/cluster/list-alerts-input.json#/properties/categories")

### categories Type

`string[]`
