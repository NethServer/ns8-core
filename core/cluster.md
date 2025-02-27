# Cluster library Schema

```txt
http://schema.nethserver.org/cluster.json
```

Cluster actions validation data formats

| Abstract               | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                          |
| :--------------------- | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------- |
| Cannot be instantiated | Yes        | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [cluster.json](cluster.json "open original schema") |

## Cluster library Type

unknown ([Cluster library](cluster.md))

# Cluster library Definitions

## Definitions group redis-pwh

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster.json#/definitions/redis-pwh"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group user-attributes

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster.json#/definitions/user-attributes"}
```

| Property                       | Type      | Required | Nullable       | Defined by                                                                                                                                                                         |
| :----------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [display\_name](#display_name) | `string`  | Optional | cannot be null | [Cluster library](cluster-definitions-user-attributes-properties-display_name.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/display_name") |
| [2fa](#2fa)                    | `boolean` | Optional | cannot be null | [Cluster library](cluster-definitions-user-attributes-properties-2fa.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/2fa")                   |

### display\_name



`display_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [Cluster library](cluster-definitions-user-attributes-properties-display_name.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/display_name")

#### display\_name Type

`string`

#### display\_name Constraints

**minimum length**: the minimum number of characters for this string is: `0`

### 2fa



`2fa`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [Cluster library](cluster-definitions-user-attributes-properties-2fa.md "http://schema.nethserver.org/cluster.json#/definitions/user-attributes/properties/2fa")

#### 2fa Type

`boolean`

## Definitions group grant-object

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster.json#/definitions/grant-object"}
```

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                                   |
| :------------ | :------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [role](#role) | `string` | Optional | cannot be null | [Cluster library](cluster-definitions-grant-object-properties-role.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/role") |
| [on](#on)     | `string` | Optional | cannot be null | [Cluster library](cluster-definitions-grant-object-properties-on.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/on")     |

### role



`role`

* is optional

* Type: `string`

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-object-properties-role.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/role")

#### role Type

`string`

#### role Constraints

**minimum length**: the minimum number of characters for this string is: `1`

### on



`on`

* is optional

* Type: `string`

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-object-properties-on.md "http://schema.nethserver.org/cluster.json#/definitions/grant-object/properties/on")

#### on Type

`string`

#### on Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## Definitions group strict-username-string

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster.json#/definitions/strict-username-string"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |

## Definitions group grant-assertion

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster.json#/definitions/grant-assertion"}
```

| Property          | Type     | Required | Nullable       | Defined by                                                                                                                                                                          |
| :---------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [action](#action) | `string` | Required | cannot be null | [Cluster library](cluster-definitions-grant-assertion-properties-action-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/action") |
| [on](#on-1)       | `string` | Required | cannot be null | [Cluster library](cluster-definitions-grant-assertion-properties-on-clause-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/on")       |
| [to](#to)         | `string` | Required | cannot be null | [Cluster library](cluster-definitions-grant-assertion-properties-user-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/to")       |

### action

Match action names with wildcard `*` character

`action`

* is required

* Type: `string` ([Action name matcher](cluster-definitions-grant-assertion-properties-action-name-matcher.md))

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-assertion-properties-action-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/action")

#### action Type

`string` ([Action name matcher](cluster-definitions-grant-assertion-properties-action-name-matcher.md))

#### action Constraints

**minimum length**: the minimum number of characters for this string is: `1`

#### action Examples

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

### on

Match the `agent_id`: a prefixed module name, node ID, or `cluster`

`on`

* is required

* Type: `string` ([On clause matcher](cluster-definitions-grant-assertion-properties-on-clause-matcher.md))

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-assertion-properties-on-clause-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/on")

#### on Type

`string` ([On clause matcher](cluster-definitions-grant-assertion-properties-on-clause-matcher.md))

#### on Constraints

**minimum length**: the minimum number of characters for this string is: `1`

#### on Examples

```json
"module/dokuwiki*"
```

```json
"node/*"
```

```json
"cluster"
```

### to

Match Redis user names with wildcard `*` character

`to`

* is required

* Type: `string` ([User name matcher](cluster-definitions-grant-assertion-properties-user-name-matcher.md))

* cannot be null

* defined in: [Cluster library](cluster-definitions-grant-assertion-properties-user-name-matcher.md "http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/to")

#### to Type

`string` ([User name matcher](cluster-definitions-grant-assertion-properties-user-name-matcher.md))

#### to Constraints

**minimum length**: the minimum number of characters for this string is: `1`

#### to Examples

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

## Definitions group ipv4-cidr

Reference this group by using

```json
{"$ref":"http://schema.nethserver.org/cluster.json#/definitions/ipv4-cidr"}
```

| Property | Type | Required | Nullable | Defined by |
| :------- | :--- | :------- | :------- | :--------- |
