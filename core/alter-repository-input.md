# alter-repository input Schema

```txt
http://schema.nethserver.org/cluster/alter-repository-input.json
```

Input schema of the alter-repository action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [alter-repository-input.json](cluster/alter-repository-input.json "open original schema") |

## alter-repository input Type

`object` ([alter-repository input](alter-repository-input.md))

## alter-repository input Examples

```json
{
  "name": "repository1",
  "testing": true,
  "status": true
}
```

# alter-repository input Properties

| Property            | Type      | Required | Nullable       | Defined by                                                                                                                                                    |
| :------------------ | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [name](#name)       | `string`  | Required | cannot be null | [alter-repository input](alter-repository-input-properties-name.md "http://schema.nethserver.org/cluster/alter-repository-input.json#/properties/name")       |
| [testing](#testing) | `boolean` | Optional | cannot be null | [alter-repository input](alter-repository-input-properties-testing.md "http://schema.nethserver.org/cluster/alter-repository-input.json#/properties/testing") |
| [status](#status)   | `boolean` | Required | cannot be null | [alter-repository input](alter-repository-input-properties-status.md "http://schema.nethserver.org/cluster/alter-repository-input.json#/properties/status")   |

## name

Unique repository name

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [alter-repository input](alter-repository-input-properties-name.md "http://schema.nethserver.org/cluster/alter-repository-input.json#/properties/name")

### name Type

`string`

## testing

Use testing releases to install new instances and update existing ones.

`testing`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [alter-repository input](alter-repository-input-properties-testing.md "http://schema.nethserver.org/cluster/alter-repository-input.json#/properties/testing")

### testing Type

`boolean`

## status

Enable or disable the repository. When 'true', the repository is enabled.

`status`

* is required

* Type: `boolean`

* cannot be null

* defined in: [alter-repository input](alter-repository-input-properties-status.md "http://schema.nethserver.org/cluster/alter-repository-input.json#/properties/status")

### status Type

`boolean`
