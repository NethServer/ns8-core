# User name matcher Schema

```txt
http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/to
```

Match Redis user names with wildcard `*` character

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [cluster.json\*](cluster.json "open original schema") |

## to Type

`string` ([User name matcher](cluster-definitions-grant-assertion-properties-user-name-matcher.md))

## to Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## to Examples

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
