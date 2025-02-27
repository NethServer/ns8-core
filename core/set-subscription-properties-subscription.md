# Untitled undefined type in set-subscription Schema

```txt
http://schema.nethserver.org/cluster/set-subscription.json#/properties/subscription
```

Setting a null subscription, clears all existing subscription attributes and stops any running service

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-subscription.json\*](cluster/set-subscription.json "open original schema") |

## subscription Type

`object` ([Details](set-subscription-properties-subscription.md))

# subscription Properties

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                                                                        |
| :------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [auth\_token](#auth_token) | `string` | Required | cannot be null | [set-subscription](set-subscription-properties-subscription-properties-auth_token.md "http://schema.nethserver.org/cluster/set-subscription.json#/properties/subscription/properties/auth_token") |

## auth\_token



`auth_token`

* is required

* Type: `string`

* cannot be null

* defined in: [set-subscription](set-subscription-properties-subscription-properties-auth_token.md "http://schema.nethserver.org/cluster/set-subscription.json#/properties/subscription/properties/auth_token")

### auth\_token Type

`string`

### auth\_token Constraints

**maximum length**: the maximum number of characters for this string is: `128`

**minimum length**: the minimum number of characters for this string is: `32`
