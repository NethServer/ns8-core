# bind-user-domains-input Schema

```txt
http://schema.nethserver.org/cluster/bind-user-domains-input.json
```

Input schema of the bind-user-domains action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [bind-user-domains-input.json](cluster/bind-user-domains-input.json "open original schema") |

## bind-user-domains-input Type

`object` ([bind-user-domains-input](bind-user-domains-input.md))

## bind-user-domains-input Examples

```json
{
  "domains": [
    "mydom.test"
  ]
}
```

# bind-user-domains-input Properties

| Property            | Type    | Required | Nullable       | Defined by                                                                                                                                                       |
| :------------------ | :------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domains](#domains) | `array` | Required | cannot be null | [bind-user-domains-input](bind-user-domains-input-properties-domains.md "http://schema.nethserver.org/cluster/bind-user-domains-input.json#/properties/domains") |

## domains

One or more domains to bind with the module calling this action

`domains`

* is required

* Type: `string[]`

* cannot be null

* defined in: [bind-user-domains-input](bind-user-domains-input-properties-domains.md "http://schema.nethserver.org/cluster/bind-user-domains-input.json#/properties/domains")

### domains Type

`string[]`
