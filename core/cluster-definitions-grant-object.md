# Grant object Schema

```txt
http://schema.nethserver.org/cluster/alter-user-input.json#/properties/revoke/items
```

A grant object establishes a relation between a role and the cluster objects matching the "on" clause

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-user-input.json\*](cluster/alter-user-input.json "open original schema") |

## items Type

`object` ([Grant object](cluster-definitions-grant-object.md))

## items Examples

```json
{
  "role": "owner",
  "on": "*"
}
```

# items Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                   |
| :------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [role](#role) | `string` | Optional | cannot be null | [Cluster library](cluster-definitions-grant-object-properties-role.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/role") |
| [on](#on)     | `string` | Optional | cannot be null | [Cluster library](cluster-definitions-grant-object-properties-on.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/on")     |

## role



`role`

* is optional

* Type: `string`

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-object-properties-role.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/role")

### role Type

`string`

### role Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## on



`on`

* is optional

* Type: `string`

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-object-properties-on.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/on")

### on Type

`string`

### on Constraints

**minimum length**: the minimum number of characters for this string is: `1`
