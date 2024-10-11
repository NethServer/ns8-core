# remove-internal-provider input Schema

```txt
http://schema.nethserver.org/cluster/remove-internal-provider-input.json
```

Safely remove a user domain provider.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-internal-provider-input.json](cluster/remove-internal-provider-input.json "open original schema") |

## remove-internal-provider input Type

`object` ([remove-internal-provider input](remove-internal-provider-input.md))

## remove-internal-provider input Examples

```json
{
  "module_id": "openldap1"
}
```

```json
{
  "module_id": "samba1",
  "adminuser": "administrator",
  "adminpass": "Nethesis,1234"
}
```

# remove-internal-provider input Properties

| Property                 | Type     | Required | Nullable       | Defined by                                                                                                                                                                                        |
| :----------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [module\_id](#module_id) | `string` | Required | cannot be null | [remove-internal-provider input](remove-internal-provider-input-properties-module-identifier.md "http://schema.nethserver.org/cluster/remove-internal-provider-input.json#/properties/module_id") |
| [adminuser](#adminuser)  | `string` | Optional | cannot be null | [remove-internal-provider input](remove-internal-provider-input-properties-adminuser.md "http://schema.nethserver.org/cluster/remove-internal-provider-input.json#/properties/adminuser")         |
| [adminpass](#adminpass)  | `string` | Optional | cannot be null | [remove-internal-provider input](remove-internal-provider-input-properties-adminpass.md "http://schema.nethserver.org/cluster/remove-internal-provider-input.json#/properties/adminpass")         |

## module\_id



`module_id`

* is required

* Type: `string` ([Module identifier](remove-internal-provider-input-properties-module-identifier.md))

* cannot be null

* defined in: [remove-internal-provider input](remove-internal-provider-input-properties-module-identifier.md "http://schema.nethserver.org/cluster/remove-internal-provider-input.json#/properties/module_id")

### module\_id Type

`string` ([Module identifier](remove-internal-provider-input-properties-module-identifier.md))

### module\_id Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## adminuser



`adminuser`

* is optional

* Type: `string`

* cannot be null

* defined in: [remove-internal-provider input](remove-internal-provider-input-properties-adminuser.md "http://schema.nethserver.org/cluster/remove-internal-provider-input.json#/properties/adminuser")

### adminuser Type

`string`

### adminuser Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## adminpass



`adminpass`

* is optional

* Type: `string`

* cannot be null

* defined in: [remove-internal-provider input](remove-internal-provider-input-properties-adminpass.md "http://schema.nethserver.org/cluster/remove-internal-provider-input.json#/properties/adminpass")

### adminpass Type

`string`

### adminpass Constraints

**minimum length**: the minimum number of characters for this string is: `1`
