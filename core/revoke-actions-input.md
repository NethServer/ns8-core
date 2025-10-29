# revoke-actions input Schema

```txt
http://schema.nethserver.org/cluster/revoke-actions-input.json
```

Revoke permissions with matches. See also the grant-actions input schema.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [revoke-actions-input.json](cluster/revoke-actions-input.json "open original schema") |

## revoke-actions input Type

`object` ([revoke-actions input](revoke-actions-input.md))

## revoke-actions input Examples

```json
{
  "on": "cluster",
  "to": "reader",
  "action": "list-something"
}
```

```json
{
  "action": "list-*",
  "to": "reader",
  "on": "cluster"
}
```

# revoke-actions input Properties

| Property          | Type          | Required | Nullable       | Defined by                                                                                                                                            |
| :---------------- | :------------ | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
| [action](#action) | Not specified | Required | cannot be null | [revoke-actions input](revoke-actions-input-properties-action.md "http://schema.nethserver.org/cluster/revoke-actions-input.json#/properties/action") |
| [to](#to)         | Not specified | Required | cannot be null | [revoke-actions input](revoke-actions-input-properties-to.md "http://schema.nethserver.org/cluster/revoke-actions-input.json#/properties/to")         |
| [on](#on)         | Not specified | Required | cannot be null | [revoke-actions input](revoke-actions-input-properties-on.md "http://schema.nethserver.org/cluster/revoke-actions-input.json#/properties/on")         |

## action



`action`

* is required

* Type: unknown

* cannot be null

* defined in: [revoke-actions input](revoke-actions-input-properties-action.md "http://schema.nethserver.org/cluster/revoke-actions-input.json#/properties/action")

### action Type

unknown

## to



`to`

* is required

* Type: unknown

* cannot be null

* defined in: [revoke-actions input](revoke-actions-input-properties-to.md "http://schema.nethserver.org/cluster/revoke-actions-input.json#/properties/to")

### to Type

unknown

## on



`on`

* is required

* Type: unknown

* cannot be null

* defined in: [revoke-actions input](revoke-actions-input-properties-on.md "http://schema.nethserver.org/cluster/revoke-actions-input.json#/properties/on")

### on Type

unknown
