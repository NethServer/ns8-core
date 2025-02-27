# remove-tun input Schema

```txt
http://schema.nethserver.org/node/remove-tun-input.json
```

Remove a tun interface

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                 |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-tun-input.json](node/remove-tun-input.json "open original schema") |

## remove-tun input Type

`object` ([remove-tun input](remove-tun-input.md))

## remove-tun input Examples

```json
{
  "tun": "tun0"
}
```

# remove-tun input Properties

| Property    | Type     | Required | Nullable       | Defined by                                                                                                                            |
| :---------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------ |
| [tun](#tun) | `string` | Required | cannot be null | [remove-tun input](remove-tun-input-properties-tun-name.md "http://schema.nethserver.org/node/remove-tun-input.json#/properties/tun") |

## tun



`tun`

* is required

* Type: `string` ([TUN name](remove-tun-input-properties-tun-name.md))

* cannot be null

* defined in: [remove-tun input](remove-tun-input-properties-tun-name.md "http://schema.nethserver.org/node/remove-tun-input.json#/properties/tun")

### tun Type

`string` ([TUN name](remove-tun-input-properties-tun-name.md))

### tun Constraints

**minimum length**: the minimum number of characters for this string is: `4`
