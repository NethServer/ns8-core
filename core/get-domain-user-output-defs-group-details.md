# Group details Schema

```txt
http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/group
```



| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                  |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-domain-user-output.json\*](cluster/get-domain-user-output.json "open original schema") |

## group Type

`object` ([Group details](get-domain-user-output-defs-group-details.md))

# group Properties

| Property                    | Type     | Required | Nullable       | Defined by                                                                                                                                                                                           |
| :-------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [group](#group)             | `string` | Required | cannot be null | [get-domain-user output](get-domain-user-output-defs-group-details-properties-group.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/group/properties/group")             |
| [description](#description) | `string` | Required | cannot be null | [get-domain-user output](get-domain-user-output-defs-group-details-properties-description.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/group/properties/description") |

## group



`group`

* is required

* Type: `string`

* cannot be null

* defined in: [get-domain-user output](get-domain-user-output-defs-group-details-properties-group.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/group/properties/group")

### group Type

`string`

### group Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## description

A brief description of the group purpose

`description`

* is required

* Type: `string` ([Description](get-domain-user-output-defs-group-details-properties-description.md))

* cannot be null

* defined in: [get-domain-user output](get-domain-user-output-defs-group-details-properties-description.md "http://schema.nethserver.org/cluster/get-domain-user-output.json#/$defs/group/properties/description")

### description Type

`string` ([Description](get-domain-user-output-defs-group-details-properties-description.md))
