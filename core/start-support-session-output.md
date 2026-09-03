# start-support-session output Schema

```txt
http://schema.nethserver.org/node/start-support-session-output.json
```

Start the support session and obtain the session ID. See get-support-session for details

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                         |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [start-support-session-output.json](node/start-support-session-output.json "open original schema") |

## start-support-session output Type

`object` ([start-support-session output](start-support-session-output.md))

## start-support-session output Examples

```json
{
  "session_id": "811630de-67f4-5b6d-8b2f-7cc01f592b1b"
}
```

# start-support-session output Properties

| Property                   | Type     | Required | Nullable       | Defined by                                                                                                                                                                         |
| :------------------------- | :------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [session\_id](#session_id) | `string` | Required | cannot be null | [start-support-session output](start-support-session-output-properties-session_id.md "http://schema.nethserver.org/node/start-support-session-output.json#/properties/session_id") |

## session\_id



`session_id`

* is required

* Type: `string`

* cannot be null

* defined in: [start-support-session output](start-support-session-output-properties-session_id.md "http://schema.nethserver.org/node/start-support-session-output.json#/properties/session_id")

### session\_id Type

`string`
