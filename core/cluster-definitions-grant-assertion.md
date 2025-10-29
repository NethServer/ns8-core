# Untitled object in grant-actions input Schema

```txt
http://schema.nethserver.org/cluster/grant-actions-input.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [grant-actions-input.json\*](cluster/grant-actions-input.json "open original schema") |

## items Type

`object` ([Details](cluster-definitions-grant-assertion.md))

# items Properties

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                                          |
| :---------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [action](#action) | `string` | Required | cannot be null | [Cluster library](cluster-definitions-grant-assertion-properties-action-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/action") |
| [on](#on)         | `string` | Required | cannot be null | [Cluster library](cluster-definitions-grant-assertion-properties-on-clause-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/on")       |
| [to](#to)         | `string` | Required | cannot be null | [Cluster library](cluster-definitions-grant-assertion-properties-user-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/to")       |

## action

Match action names with wildcard `*` character

`action`

* is required

* Type: `string` ([Action name matcher](cluster-definitions-grant-assertion-properties-action-name-matcher.md))

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-assertion-properties-action-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/action")

### action Type

`string` ([Action name matcher](cluster-definitions-grant-assertion-properties-action-name-matcher.md))

### action Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### action Examples

```json
"list-*"
```

```json
"show-*"
```

```json
"list-users"
```

```json
"add-user"
```

## on

Match the `agent_id`: a prefixed module name, node ID, or `cluster`

`on`

* is required

* Type: `string` ([On clause matcher](cluster-definitions-grant-assertion-properties-on-clause-matcher.md))

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-assertion-properties-on-clause-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/on")

### on Type

`string` ([On clause matcher](cluster-definitions-grant-assertion-properties-on-clause-matcher.md))

### on Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### on Examples

```json
"module/dokuwiki*"
```

```json
"node/*"
```

```json
"cluster"
```

## to

Match Redis user names with wildcard `*` character

`to`

* is required

* Type: `string` ([User name matcher](cluster-definitions-grant-assertion-properties-user-name-matcher.md))

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-assertion-properties-user-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/to")

### to Type

`string` ([User name matcher](cluster-definitions-grant-assertion-properties-user-name-matcher.md))

### to Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### to Examples

```json
"admin"
```

```json
"cluster"
```

```json
"module/*"
```

```json
"node/3"
```

```json
"*"
```
