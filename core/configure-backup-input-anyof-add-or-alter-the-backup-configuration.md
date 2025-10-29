# Add or alter the backup configuration Schema

```txt
http://schema.nethserver.org/agent/configure-backup-input.json#/anyOf/0
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [configure-backup-input.json\*](agent/configure-backup-input.json "open original schema") |

## 0 Type

unknown ([Add or alter the backup configuration](configure-backup-input-anyof-add-or-alter-the-backup-configuration.md))

# 0 Properties

| Property          | Type          | Required | Nullable       | Defined by                                                                                                                                                                                                    |
| :---------------- | :------------ | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [optype](#optype) | Not specified | Required | cannot be null | [configure-backup input](configure-backup-input-anyof-add-or-alter-the-backup-configuration-properties-optype.md "http://schema.nethserver.org/agent/configure-backup-input.json#/anyOf/0/properties/optype") |

## optype



`optype`

* is required

* Type: unknown

* cannot be null

* defined in: [configure-backup input](configure-backup-input-anyof-add-or-alter-the-backup-configuration-properties-optype.md "http://schema.nethserver.org/agent/configure-backup-input.json#/anyOf/0/properties/optype")

### optype Type

unknown

### optype Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value     | Explanation |
| :-------- | :---------- |
| `"add"`   |             |
| `"alter"` |             |
