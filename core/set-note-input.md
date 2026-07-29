# set-note input Schema

```txt
http://schema.nethserver.org/agent/set-note-input.json
```

A short note to help identify or describe a module instance. The note is visible in Applications Center.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-note-input.json](agent/set-note-input.json "open original schema") |

## set-note input Type

`object` ([set-note input](set-note-input.md))

## set-note input Examples

```json
{
  "note": "this module is for my personal use"
}
```

# set-note input Properties

| Property      | Type     | Required | Nullable       | Defined by                                                                                                                    |
| :------------ | :------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------- |
| [note](#note) | `string` | Required | cannot be null | [set-note input](set-note-input-properties-note.md "http://schema.nethserver.org/agent/set-note-input.json#/properties/note") |

## note



`note`

* is required

* Type: `string`

* cannot be null

* defined in: [set-note input](set-note-input-properties-note.md "http://schema.nethserver.org/agent/set-note-input.json#/properties/note")

### note Type

`string`

### note Constraints

**maximum length**: the maximum number of characters for this string is: `100`
