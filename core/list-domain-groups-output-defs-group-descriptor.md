# Group descriptor Schema

```txt
http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group
```

Basic description of a group: name and description

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [list-domain-groups-output.json\*](cluster/list-domain-groups-output.json "open original schema") |

## group Type

`object` ([Group descriptor](list-domain-groups-output-defs-group-descriptor.md))

# group Properties

| Property                    | Type     | Required | Nullable       | Defined by                                                                                                                                                                                                       |
| :-------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [group](#group)             | `string` | Required | cannot be null | [list-domain-groups output](list-domain-groups-output-defs-group-descriptor-properties-group.md "http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group/properties/group")             |
| [description](#description) | `string` | Required | cannot be null | [list-domain-groups output](list-domain-groups-output-defs-group-descriptor-properties-description.md "http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group/properties/description") |
| [users](#users)             | `array`  | Optional | cannot be null | [list-domain-groups output](list-domain-groups-output-defs-group-descriptor-properties-group-members.md "http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group/properties/users")     |

## group



`group`

* is required

* Type: `string`

* cannot be null

* defined in: [list-domain-groups output](list-domain-groups-output-defs-group-descriptor-properties-group.md "http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group/properties/group")

### group Type

`string`

### group Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## description

A brief description of the group purpose

`description`

* is required

* Type: `string` ([Description](list-domain-groups-output-defs-group-descriptor-properties-description.md))

* cannot be null

* defined in: [list-domain-groups output](list-domain-groups-output-defs-group-descriptor-properties-description.md "http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group/properties/description")

### description Type

`string` ([Description](list-domain-groups-output-defs-group-descriptor-properties-description.md))

## users



`users`

* is optional

* Type: `string[]` ([Member name](list-domain-groups-output-defs-group-descriptor-properties-group-members-member-name.md))

* cannot be null

* defined in: [list-domain-groups output](list-domain-groups-output-defs-group-descriptor-properties-group-members.md "http://schema.nethserver.org/cluster/list-domain-groups-output.json#/$defs/group/properties/users")

### users Type

`string[]` ([Member name](list-domain-groups-output-defs-group-descriptor-properties-group-members-member-name.md))
