# Counters Schema

```txt
http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters
```

The cached number of users and groups returned by their respective last API calls

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-user-domains-output.json\*](cluster/list-user-domains-output.json "open original schema") |

## counters Type

`object` ([Counters](list-user-domains-output-defs-user-domain-properties-counters.md))

# counters Properties

| Property          | Type      | Required | Nullable    | Defined by                                                                                                                                                                                                                                          |
| :---------------- | :-------- | :------- | :---------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [users](#users)   | `integer` | Required | can be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters-properties-user-counter.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters/properties/users")   |
| [groups](#groups) | `integer` | Required | can be null | [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters-properties-group-counter.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters/properties/groups") |

## users



`users`

* is required

* Type: `integer` ([User counter](list-user-domains-output-defs-user-domain-properties-counters-properties-user-counter.md))

* can be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters-properties-user-counter.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters/properties/users")

### users Type

`integer` ([User counter](list-user-domains-output-defs-user-domain-properties-counters-properties-user-counter.md))

## groups



`groups`

* is required

* Type: `integer` ([Group counter](list-user-domains-output-defs-user-domain-properties-counters-properties-group-counter.md))

* can be null

* defined in: [list-user-domains output](list-user-domains-output-defs-user-domain-properties-counters-properties-group-counter.md "http://schema.nethserver.org/cluster/list-user-domains-output.json#/$defs/user-domain/properties/counters/properties/groups")

### groups Type

`integer` ([Group counter](list-user-domains-output-defs-user-domain-properties-counters-properties-group-counter.md))
