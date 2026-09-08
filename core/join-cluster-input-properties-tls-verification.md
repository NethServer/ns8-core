# TLS verification Schema

```txt
http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/tls_verify
```

If set to `false` the given @url is not subject to TLS certificate validation

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [join-cluster-input.json\*](cluster/join-cluster-input.json "open original schema") |

## tls\_verify Type

`boolean` ([TLS verification](join-cluster-input-properties-tls-verification.md))

## tls\_verify Examples

```json
true
```

```json
false
```
