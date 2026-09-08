# add-module output Schema

```txt
http://schema.nethserver.org/node/add-module-output.json
```

Return generated information of the newly added module

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                   |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-module-output.json](node/add-module-output.json "open original schema") |

## add-module output Type

`object` ([add-module output](add-module-output.md))

## add-module output Examples

```json
{
  "redis_sha256": "a0c85889af15c49d7dd08cb896ec0c5e960a67d5d48d51f06757c2cee3a0277b"
}
```

# add-module output Properties

| Property                       | Type     | Required | Nullable       | Defined by                                                                                                                                            |
| :----------------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------- |
| [redis\_sha256](#redis_sha256) | `string` | Required | cannot be null | [add-module output](add-module-output-properties-redis_sha256.md "http://schema.nethserver.org/node/add-module-output.json#/properties/redis_sha256") |

## redis\_sha256

The generated Redis hashed password for the new node

`redis_sha256`

* is required

* Type: `string`

* cannot be null

* defined in: [add-module output](add-module-output-properties-redis_sha256.md "http://schema.nethserver.org/node/add-module-output.json#/properties/redis_sha256")

### redis\_sha256 Type

`string`
