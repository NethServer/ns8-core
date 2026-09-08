# restart-module input Schema

```txt
http://schema.nethserver.org/node/restart-module-input.json
```

This schema validates the input for restarting a module's services.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                         |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [restart-module-input.json](node/restart-module-input.json "open original schema") |

## restart-module input Type

`object` ([restart-module input](restart-module-input.md))

## restart-module input Examples

```json
{
  "module_id": "mymodule2"
}
```

# restart-module input Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                       |
| :----------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id) | `string` | Required | cannot be null | [restart-module input](restart-module-input-properties-module-identifier.md "http://schema.nethserver.org/node/restart-module-input.json#/properties/module_id") |

## module\_id

The unique identifier of the module whose services are to be restarted.

`module_id`

* is required

* Type: `string` ([Module identifier](restart-module-input-properties-module-identifier.md))

* cannot be null

* defined in: [restart-module input](restart-module-input-properties-module-identifier.md "http://schema.nethserver.org/node/restart-module-input.json#/properties/module_id")

### module\_id Type

`string` ([Module identifier](restart-module-input-properties-module-identifier.md))

### module\_id Constraints

**minimum length**: the minimum number of characters for this string is: `1`
