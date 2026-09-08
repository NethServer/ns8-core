# Action name matcher Schema

```txt
http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/action
```

Match action names with wildcard `*` character

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [cluster.json\*](cluster.json "open original schema") |

## action Type

`string` ([Action name matcher](cluster-definitions-grant-assertion-properties-action-name-matcher.md))

## action Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## action Examples

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
