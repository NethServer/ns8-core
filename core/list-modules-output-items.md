# Untitled object in list-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-modules-output.json#/items
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                            |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-modules-output.json\*](cluster/list-modules-output.json "open original schema") |

## items Type

`object` ([Details](list-modules-output-items.md))

# items Properties

| Property                                   | Type     | Required | Nullable       | Defined by                                                                                                                                                                             |
| :----------------------------------------- | :------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)                              | `string` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/name")                             |
| [description](#description)                | `object` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-description.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/description")               |
| [logo](#logo)                              | `string` | Required | can be null    | [list-modules output](list-modules-output-items-properties-logo.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/logo")                             |
| [categories](#categories)                  | `array`  | Required | cannot be null | [list-modules output](list-modules-output-items-properties-categories.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/categories")                 |
| [authors](#authors)                        | `array`  | Required | cannot be null | [list-modules output](list-modules-output-items-properties-authors.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/authors")                       |
| [docs](#docs)                              | `object` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-docs.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/docs")                             |
| [source](#source)                          | `string` | Required | cannot be null | [list-modules output](list-modules-output-items-properties-source.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/source")                         |
| [repository](#repository)                  | `string` | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-repository.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository")                 |
| [repository\_updated](#repository_updated) | `string` | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-repository_updated.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository_updated") |
| [versions](#versions)                      | `array`  | Required | cannot be null | [list-modules output](list-modules-output-items-properties-versions.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions")                     |

## name

Unique name of a package

`name`

*   is required

*   Type: `string`

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/name")

### name Type

`string`

## description

A map of language codes (eg. en, it) with the translated description

`description`

*   is required

*   Type: `object` ([Details](list-modules-output-items-properties-description.md))

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-description.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/description")

### description Type

`object` ([Details](list-modules-output-items-properties-description.md))

## logo

The filename of the logo. The logo must be a PNG image of 256x256 pixels

`logo`

*   is required

*   Type: `string`

*   can be null

*   defined in: [list-modules output](list-modules-output-items-properties-logo.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/logo")

### logo Type

`string`

## categories



`categories`

*   is required

*   Type: `string[]`

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-categories.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/categories")

### categories Type

`string[]`

## authors



`authors`

*   is required

*   Type: `object[]` ([Details](list-modules-output-items-properties-authors-items.md))

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-authors.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/authors")

### authors Type

`object[]` ([Details](list-modules-output-items-properties-authors-items.md))

## docs



`docs`

*   is required

*   Type: `object` ([Details](list-modules-output-items-properties-docs.md))

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-docs.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/docs")

### docs Type

`object` ([Details](list-modules-output-items-properties-docs.md))

## source

URL of package inside a container registry. The URL must not contain the prefix 'docker://' or similar

`source`

*   is required

*   Type: `string`

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-source.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/source")

### source Type

`string`

## repository

The internal ID of the repository inside redis

`repository`

*   is optional

*   Type: `string`

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-repository.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository")

### repository Type

`string`

## repository\_updated

Date and time of last modification to remote repodata

`repository_updated`

*   is optional

*   Type: `string`

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-repository_updated.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository_updated")

### repository\_updated Type

`string`

## versions



`versions`

*   is required

*   Type: `object[]` ([Details](list-modules-output-items-properties-versions-items.md))

*   cannot be null

*   defined in: [list-modules output](list-modules-output-items-properties-versions.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions")

### versions Type

`object[]` ([Details](list-modules-output-items-properties-versions-items.md))
