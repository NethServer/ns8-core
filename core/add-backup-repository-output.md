# add-backup-repository output Schema

```txt
http://schema.nethserver.org/cluster/add-backup-repository-output.json
```

Output schema of the add-backup-repository action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [add-backup-repository-output.json](cluster/add-backup-repository-output.json "open original schema") |

## add-backup-repository output Type

`object` ([add-backup-repository output](add-backup-repository-output.md))

## add-backup-repository output Examples

```json
{
  "password": "55cdf59499dc1b5c625c8e20682f7c34302790767c2a13abd20d4289db0feca6",
  "id": "ce712dfe-b1d0-47f1-a149-a14d1203aab8"
}
```

# add-backup-repository output Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                                                                                            |
| :-------------------- | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [id](#id)             | `string` | Required | cannot be null | [add-backup-repository output](add-backup-repository-output-properties-identifier-of-the-added-repository.md "http://schema.nethserver.org/cluster/add-backup-repository-output.json#/properties/id") |
| [password](#password) | `string` | Required | cannot be null | [add-backup-repository output](add-backup-repository-output-properties-generated-random-password.md "http://schema.nethserver.org/cluster/add-backup-repository-output.json#/properties/password")    |

## id



`id`

* is required

* Type: `string` ([Identifier of the added repository](add-backup-repository-output-properties-identifier-of-the-added-repository.md))

* cannot be null

* defined in: [add-backup-repository output](add-backup-repository-output-properties-identifier-of-the-added-repository.md "http://schema.nethserver.org/cluster/add-backup-repository-output.json#/properties/id")

### id Type

`string` ([Identifier of the added repository](add-backup-repository-output-properties-identifier-of-the-added-repository.md))

## password



`password`

* is required

* Type: `string` ([Generated random password](add-backup-repository-output-properties-generated-random-password.md))

* cannot be null

* defined in: [add-backup-repository output](add-backup-repository-output-properties-generated-random-password.md "http://schema.nethserver.org/cluster/add-backup-repository-output.json#/properties/password")

### password Type

`string` ([Generated random password](add-backup-repository-output-properties-generated-random-password.md))
