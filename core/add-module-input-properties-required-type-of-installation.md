# Required type of installation Schema

```txt
http://schema.nethserver.org/node/add-module-input.json#/properties/is_rootfull
```

If `true` the module is installed in rootfull mode.
Its agent runs as root and its containers can be granted **privileged access**.

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                   |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [add-module-input.json\*](node/add-module-input.json "open original schema") |

## is\_rootfull Type

`boolean` ([Required type of installation](add-module-input-properties-required-type-of-installation.md))
