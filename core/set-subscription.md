# set-subscription Schema

```txt
http://schema.nethserver.org/cluster/set-subscription.json
```

Set up support subscription

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                    |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-subscription.json](cluster/set-subscription.json "open original schema") |

## set-subscription Type

`object` ([set-subscription](set-subscription.md))

## set-subscription Examples

```json
{
  "subscription": {
    "auth_token": "d489701d76b1424e9dacb08b73741ef6d489701d76b1424e9dacb08b73741ef6"
  }
}
```

# set-subscription Properties

| Property                      | Type     | Required | Nullable    | Defined by                                                                                                                                            |
| :---------------------------- | :------- | :------- | :---------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
| [subscription](#subscription) | `object` | Required | can be null | [set-subscription](set-subscription-properties-subscription.md "http://schema.nethserver.org/cluster/set-subscription.json#/properties/subscription") |

## subscription

Setting a null subscription, clears all existing subscription attributes and stops any running service

`subscription`

* is required

* Type: `object` ([Details](set-subscription-properties-subscription.md))

* can be null

* defined in: [set-subscription](set-subscription-properties-subscription.md "http://schema.nethserver.org/cluster/set-subscription.json#/properties/subscription")

### subscription Type

`object` ([Details](set-subscription-properties-subscription.md))
