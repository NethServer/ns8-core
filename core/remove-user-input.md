# remove-user input Schema

```txt
http://schema.nethserver.org/cluster/remove-user-input.json
```

Remove a module from the cluster, optionally erasing any disk data

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                      |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [remove-user-input.json](cluster/remove-user-input.json "open original schema") |

## remove-user input Type

`object` ([remove-user input](remove-user-input.md))

## remove-user input Examples

```json
{
  "user": "jdoe"
}
```

# remove-user input Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                                        |
| :------------ | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ |
| [user](#user) | `string` | Required | cannot be null | [remove-user input](cluster-definitions-strict-username-format.md "http://schema.nethserver.org/cluster/remove-user-input.json#/properties/user") |

## user

A username string can be written in many ways. This is a quite strict form. See <https://www.unix.com/man-page/linux/8/useradd/>

`user`

* is required

* Type: `string` ([Strict username format](cluster-definitions-strict-username-format.md))

* cannot be null

* defined in: [remove-user input](cluster-definitions-strict-username-format.md "http://schema.nethserver.org/cluster/remove-user-input.json#/properties/user")

### user Type

`string` ([Strict username format](cluster-definitions-strict-username-format.md))

### user Constraints

**maximum length**: the maximum number of characters for this string is: `32`

**minimum length**: the minimum number of characters for this string is: `1`

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^[a-z_][a-z0-9_-]*[$]?$
```

[try pattern](https://regexr.com/?expression=%5E%5Ba-z_%5D%5Ba-z0-9_-%5D*%5B%24%5D%3F%24 "try regular expression with regexr.com")

### user Examples

```json
"admin"
```

```json
"u0000"
```

```json
"worker$"
```

```json
"test-user"
```

```json
"test_user"
```
