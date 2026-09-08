# set-automatic-updates Schema

```txt
http://schema.nethserver.org/cluster/set-automatic-updates-input.json
```

Configure the cluster-wide automatic updates switch and the per-instance automatic updates policy

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-automatic-updates-input.json](cluster/set-automatic-updates-input.json "open original schema") |

## set-automatic-updates Type

`object` ([set-automatic-updates](set-automatic-updates-input.md))

## set-automatic-updates Examples

```json
{
  "apply_updates_is_active": true,
  "instances": {
    "dokuwiki1": false,
    "dokuwiki2": true
  }
}
```

# set-automatic-updates Properties

| Property                                               | Type      | Required | Nullable       | Defined by                                                                                                                                                                                             |
| :----------------------------------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [apply\_updates\_is\_active](#apply_updates_is_active) | `boolean` | Optional | cannot be null | [set-automatic-updates](set-automatic-updates-input-properties-apply_updates_is_active.md "http://schema.nethserver.org/cluster/set-automatic-updates-input.json#/properties/apply_updates_is_active") |
| [instances](#instances)                                | `object`  | Optional | cannot be null | [set-automatic-updates](set-automatic-updates-input-properties-instances.md "http://schema.nethserver.org/cluster/set-automatic-updates-input.json#/properties/instances")                             |

## apply\_updates\_is\_active

Enable or disable the cluster-wide automatic updates timer

`apply_updates_is_active`

* is optional

* Type: `boolean`

* cannot be null

* defined in: [set-automatic-updates](set-automatic-updates-input-properties-apply_updates_is_active.md "http://schema.nethserver.org/cluster/set-automatic-updates-input.json#/properties/apply_updates_is_active")

### apply\_updates\_is\_active Type

`boolean`

## instances

Map of module instance IDs to their automatic updates policy

`instances`

* is optional

* Type: `object` ([Details](set-automatic-updates-input-properties-instances.md))

* cannot be null

* defined in: [set-automatic-updates](set-automatic-updates-input-properties-instances.md "http://schema.nethserver.org/cluster/set-automatic-updates-input.json#/properties/instances")

### instances Type

`object` ([Details](set-automatic-updates-input-properties-instances.md))
