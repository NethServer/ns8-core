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

| Property                                            | Type          | Required | Nullable       | Defined by                                                                                                                                                                                       |
| :-------------------------------------------------- | :------------ | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [name](#name)                                       | `string`      | Required | cannot be null | [list-modules output](list-modules-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/name")                                       |
| [upstream_name](#upstream_name)                     | `string`      | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-upstream_name.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/upstream_name")                     |
| [certification_level](#certification_level)         | `integer`     | Required | cannot be null | [list-modules output](list-modules-output-items-properties-certification_level.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/certification_level")         |
| [rootfull](#rootfull)                               | `boolean`     | Required | cannot be null | [list-modules output](list-modules-output-items-properties-rootfull.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/rootfull")                               |
| [description](#description)                         | `object`      | Required | cannot be null | [list-modules output](list-modules-output-items-properties-description.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/description")                         |
| [logo](#logo)                                       | `string`      | Required | can be null    | [list-modules output](list-modules-output-items-properties-logo.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/logo")                                       |
| [categories](#categories)                           | `array`       | Required | cannot be null | [list-modules output](list-modules-output-items-properties-categories.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/categories")                           |
| [authors](#authors)                                 | `array`       | Required | cannot be null | [list-modules output](list-modules-output-items-properties-authors.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/authors")                                 |
| [docs](#docs)                                       | `object`      | Required | cannot be null | [list-modules output](list-modules-output-items-properties-docs.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/docs")                                       |
| [source](#source)                                   | `string`      | Required | cannot be null | [list-modules output](list-modules-output-items-properties-source.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/source")                                   |
| [repository](#repository)                           | `string`      | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-repository.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository")                           |
| [repository_updated](#repository_updated)           | `string`      | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-repository_updated.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository_updated")           |
| [disabled_updates_reason](#disabled_updates_reason) | Not specified | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-disabled_updates_reason.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/disabled_updates_reason") |
| [apply_updates_is_active](#apply_updates_is_active) | `boolean`     | Required | cannot be null | [list-modules output](list-modules-output-items-properties-apply_updates_is_active.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/apply_updates_is_active") |
| [subscription_is_active](#subscription_is_active)   | `boolean`     | Required | cannot be null | [list-modules output](list-modules-output-items-properties-subscription_is_active.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/subscription_is_active")   |
| [updates](#updates)                                 | `array`       | Required | cannot be null | [list-modules output](list-modules-output-items-properties-updates.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/updates")                                 |
| [install_destinations](#install_destinations)       | `array`       | Required | cannot be null | [list-modules output](list-modules-output-items-properties-install_destinations.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations")       |
| [installed](#installed)                             | `array`       | Required | cannot be null | [list-modules output](list-modules-output-items-properties-installed.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/installed")                             |
| [no_version_reason](#no_version_reason)             | `object`      | Optional | cannot be null | [list-modules output](list-modules-output-items-properties-no_version_reason.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason")             |
| [versions](#versions)                               | `array`       | Required | cannot be null | [list-modules output](list-modules-output-items-properties-versions.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions")                               |

## name

Unique name of a package

`name`

* is required

* Type: `string`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-name.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/name")

### name Type

`string`

## upstream_name

The alternative software name and version number, if they differ from the package name and version

`upstream_name`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-upstream_name.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/upstream_name")

### upstream_name Type

`string`

## certification_level

The higher, the better: 0=unknown certification, 5=max

`certification_level`

* is required

* Type: `integer`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-certification_level.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/certification_level")

### certification_level Type

`integer`

## rootfull

True if the application gains full OS privileges when installed

`rootfull`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-rootfull.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/rootfull")

### rootfull Type

`boolean`

## description

A map of language codes (eg. en, it) with the translated description

`description`

* is required

* Type: `object` ([Details](list-modules-output-items-properties-description.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-description.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/description")

### description Type

`object` ([Details](list-modules-output-items-properties-description.md))

## logo

The filename of the logo. The logo must be a PNG image of 256x256 pixels

`logo`

* is required

* Type: `string`

* can be null

* defined in: [list-modules output](list-modules-output-items-properties-logo.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/logo")

### logo Type

`string`

## categories



`categories`

* is required

* Type: `string[]`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-categories.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/categories")

### categories Type

`string[]`

## authors



`authors`

* is required

* Type: `object[]` ([Details](list-modules-output-items-properties-authors-items.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-authors.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/authors")

### authors Type

`object[]` ([Details](list-modules-output-items-properties-authors-items.md))

## docs



`docs`

* is required

* Type: `object` ([Details](list-modules-output-items-properties-docs.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-docs.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/docs")

### docs Type

`object` ([Details](list-modules-output-items-properties-docs.md))

## source

URL of package inside a container registry. The URL must not contain the prefix 'docker://' or similar

`source`

* is required

* Type: `string`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-source.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/source")

### source Type

`string`

## repository

The internal ID of the repository inside redis

`repository`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-repository.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository")

### repository Type

`string`

## repository_updated

Date and time of last modification to remote repodata

`repository_updated`

* is optional

* Type: `string`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-repository_updated.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/repository_updated")

### repository_updated Type

`string`

## disabled_updates_reason

A non-empty strings indicates that updates are disabled for some reason. The string value identifies the reason.

`disabled_updates_reason`

* is optional

* Type: unknown

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-disabled_updates_reason.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/disabled_updates_reason")

### disabled_updates_reason Type

unknown

### disabled_updates_reason Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value             | Explanation |
| :---------------- | :---------- |
| `""`              |             |
| `"ns7_migration"` |             |

## apply_updates_is_active

True if the cluster-wide automatic updates timer is active

`apply_updates_is_active`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-apply_updates_is_active.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/apply_updates_is_active")

### apply_updates_is_active Type

`boolean`

## subscription_is_active

True if a Nethesis subscription is active on the cluster

`subscription_is_active`

* is required

* Type: `boolean`

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-subscription_is_active.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/subscription_is_active")

### subscription_is_active Type

`boolean`

## updates



`updates`

* is required

* Type: `object[]` ([Details](list-modules-output-items-properties-updates-items.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-updates.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/updates")

### updates Type

`object[]` ([Details](list-modules-output-items-properties-updates-items.md))

## install_destinations

Describe for each node of the cluster if the node is eligible or not to install a new module instance. If not, a reject reason is returned.

`install_destinations`

* is required

* Type: `object[]` ([Details](list-modules-output-items-properties-install_destinations-items.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-install_destinations.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/install_destinations")

### install_destinations Type

`object[]` ([Details](list-modules-output-items-properties-install_destinations-items.md))

## installed



`installed`

* is required

* Type: `object[]` ([Details](list-modules-output-items-properties-installed-items.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-installed.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/installed")

### installed Type

`object[]` ([Details](list-modules-output-items-properties-installed-items.md))

## no_version_reason

If the versions array is empty, this object is present and explains why.

`no_version_reason`

* is optional

* Type: `object` ([Details](list-modules-output-items-properties-no_version_reason.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-no_version_reason.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/no_version_reason")

### no_version_reason Type

`object` ([Details](list-modules-output-items-properties-no_version_reason.md))

## versions



`versions`

* is required

* Type: `object[]` ([Details](list-modules-output-items-properties-versions-items.md))

* cannot be null

* defined in: [list-modules output](list-modules-output-items-properties-versions.md "http://schema.nethserver.org/cluster/list-modules-output.json#/items/properties/versions")

### versions Type

`object[]` ([Details](list-modules-output-items-properties-versions-items.md))
