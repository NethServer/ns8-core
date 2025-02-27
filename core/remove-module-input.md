# remove-module input Schema

```txt
http://schema.nethserver.org/node/remove-module-input.json
```

Remove a module from the running node, optionally erasing any disk data

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                       |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-module-input.json](node/remove-module-input.json "open original schema") |

## remove-module input Type

`object` ([remove-module input](remove-module-input.md))

## remove-module input Examples

```json
{
  "module_id": "dokuwiki1",
  "preserve_data": false
}
```

# remove-module input Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                         |
| :------------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id)         | `string`  | Required | cannot be null | [remove-module input](remove-module-input-properties-module-identifier.md "http://schema.nethserver.org/node/remove-module-input.json#/properties/module_id")      |
| [preserve\_data](#preserve_data) | `boolean` | Required | cannot be null | [remove-module input](remove-module-input-properties-preserve-data-flag.md "http://schema.nethserver.org/node/remove-module-input.json#/properties/preserve_data") |

## module\_id



`module_id`

* is required

* Type: `string` ([Module identifier](remove-module-input-properties-module-identifier.md))

* cannot be null

* defined in: [remove-module input](remove-module-input-properties-module-identifier.md "http://schema.nethserver.org/node/remove-module-input.json#/properties/module_id")

### module\_id Type

`string` ([Module identifier](remove-module-input-properties-module-identifier.md))

### module\_id Examples

```json
"dokuwiki1"
```

```json
"nextcloud3"
```

## preserve\_data

If set to `false`, any module data on disk is erased

`preserve_data`

* is required

* Type: `boolean` ([Preserve data flag](remove-module-input-properties-preserve-data-flag.md))

* cannot be null

* defined in: [remove-module input](remove-module-input-properties-preserve-data-flag.md "http://schema.nethserver.org/node/remove-module-input.json#/properties/preserve_data")

### preserve\_data Type

`boolean` ([Preserve data flag](remove-module-input-properties-preserve-data-flag.md))
