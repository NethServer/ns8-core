# remove-module input Schema

```txt
http://schema.nethserver.org/cluster/remove-module-input.json
```

Remove a module from the cluster, optionally erasing any disk data

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-module-input.json](cluster/remove-module-input.json "open original schema") |

## remove-module input Type

`object` ([remove-module input](remove-module-input-1.md))

## remove-module input Examples

```json
{
  "module_id": "dokuwiki3",
  "preserve_data": false,
  "force": false
}
```

# remove-module input Properties

| Property                         | Type      | Required | Nullable       | Defined by                                                                                                                                                                                         |
| :------------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [module\_id](#module_id)         | `string`  | Required | cannot be null | [remove-module input](remove-module-input-1-properties-module-identifier.md "http://schema.nethserver.org/cluster/remove-module-input.json#/properties/module_id")                                 |
| [preserve\_data](#preserve_data) | `boolean` | Required | cannot be null | [remove-module input](remove-module-input-1-properties-preserve-data-flag.md "http://schema.nethserver.org/cluster/remove-module-input.json#/properties/preserve_data")                            |
| [force](#force)                  | `boolean` | Optional | cannot be null | [remove-module input](remove-module-input-1-properties-force-removal-from-the-cluster-even-if-not-responding.md "http://schema.nethserver.org/cluster/remove-module-input.json#/properties/force") |

## module\_id



`module_id`

* is required

* Type: `string` ([Module identifier](remove-module-input-1-properties-module-identifier.md))

* cannot be null

* defined in: [remove-module input](remove-module-input-1-properties-module-identifier.md "http://schema.nethserver.org/cluster/remove-module-input.json#/properties/module_id")

### module\_id Type

`string` ([Module identifier](remove-module-input-1-properties-module-identifier.md))

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

* Type: `boolean` ([Preserve data flag](remove-module-input-1-properties-preserve-data-flag.md))

* cannot be null

* defined in: [remove-module input](remove-module-input-1-properties-preserve-data-flag.md "http://schema.nethserver.org/cluster/remove-module-input.json#/properties/preserve_data")

### preserve\_data Type

`boolean` ([Preserve data flag](remove-module-input-1-properties-preserve-data-flag.md))

## force

The module instance is erased from Redis

`force`

* is optional

* Type: `boolean` ([Force removal from the cluster even if not responding](remove-module-input-1-properties-force-removal-from-the-cluster-even-if-not-responding.md))

* cannot be null

* defined in: [remove-module input](remove-module-input-1-properties-force-removal-from-the-cluster-even-if-not-responding.md "http://schema.nethserver.org/cluster/remove-module-input.json#/properties/force")

### force Type

`boolean` ([Force removal from the cluster even if not responding](remove-module-input-1-properties-force-removal-from-the-cluster-even-if-not-responding.md))
