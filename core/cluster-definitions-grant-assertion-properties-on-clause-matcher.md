# On clause matcher Schema

```txt
http://schema.nethserver.org/cluster.json#/definitions/grant-assertion/properties/on
```

Match the `agent_id`: a prefixed module name, node ID, or `cluster`

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                            |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [cluster.json\*](cluster.json "open original schema") |

## on Type

`string` ([On clause matcher](cluster-definitions-grant-assertion-properties-on-clause-matcher.md))

## on Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## on Examples

```json
"module/dokuwiki*"
```

```json
"node/*"
```

```json
"cluster"
```
