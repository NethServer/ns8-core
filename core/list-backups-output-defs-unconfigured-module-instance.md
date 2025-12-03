# Unconfigured module instance Schema

```txt
http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item
```

Instance with no backup configured

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-backups-output.json\*](cluster/list-backups-output.json "open original schema") |

## uinstance-item Type

`object` ([Unconfigured module instance](list-backups-output-defs-unconfigured-module-instance.md))

# uinstance-item Properties

| Property             | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                         |
| :------------------- | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)            | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/id")           |
| [ui\_name](#ui_name) | `string` | Optional | cannot be null | [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/ui_name") |

## id



`id`

* is optional

* Type: `string` ([Module ID](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/id")

### id Type

`string` ([Module ID](list-backups-output-defs-unconfigured-module-instance-properties-module-id.md))

### id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## ui\_name



`ui_name`

* is optional

* Type: `string` ([Module UI name](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md))

* cannot be null

* defined in: [list-backups output](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md "http://schema.nethserver.org/cluster/list-backups-output.json#/$defs/uinstance-item/properties/ui_name")

### ui\_name Type

`string` ([Module UI name](list-backups-output-defs-unconfigured-module-instance-properties-module-ui-name.md))
