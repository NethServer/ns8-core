# Remove the backup configuration Schema

```txt
http://schema.nethserver.org/agent/configure-backup-input.json#/anyOf/1
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [configure-backup-input.json\*](agent/configure-backup-input.json "open original schema") |

## 1 Type

unknown ([Remove the backup configuration](configure-backup-input-anyof-remove-the-backup-configuration.md))

# 1 Properties

| Property          | Type          | Required | Nullable       | Defined by                                                                                                                                                                                              |
| :---------------- | :------------ | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [optype](#optype) | Not specified | Required | cannot be null | [configure-backup input](configure-backup-input-anyof-remove-the-backup-configuration-properties-optype.md "http://schema.nethserver.org/agent/configure-backup-input.json#/anyOf/1/properties/optype") |

## optype



`optype`

* is required

* Type: unknown

* cannot be null

* defined in: [configure-backup input](configure-backup-input-anyof-remove-the-backup-configuration-properties-optype.md "http://schema.nethserver.org/agent/configure-backup-input.json#/anyOf/1/properties/optype")

### optype Type

unknown

### optype Constraints

**constant**: the value of this property must be equal to:

```json
"remove"
```
