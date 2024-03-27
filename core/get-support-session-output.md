# get-support-session output Schema

```txt
http://schema.nethserver.org/node/get-support-session-output.json
```

Return the support session status

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                     |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-support-session-output.json](node/get-support-session-output.json "open original schema") |

## get-support-session output Type

`object` ([get-support-session output](get-support-session-output.md))

## get-support-session output Examples

```json
{
  "session_id": "811630de-67f4-5b6d-8b2f-7cc01f592b1b",
  "active": true
}
```

```json
{
  "session_id": "811630de-67f4-5b6d-8b2f-7cc01f592b1b",
  "active": false
}
```

# get-support-session output Properties

| Property                   | Type      | Required | Nullable       | Defined by                                                                                                                                                                   |
| :------------------------- | :-------- | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [session\_id](#session_id) | `string`  | Required | cannot be null | [get-support-session output](get-support-session-output-properties-session_id.md "http://schema.nethserver.org/node/get-support-session-output.json#/properties/session_id") |
| [active](#active)          | `boolean` | Required | cannot be null | [get-support-session output](get-support-session-output-properties-active.md "http://schema.nethserver.org/node/get-support-session-output.json#/properties/active")         |

## session\_id

Last opened session identifier for the support team

`session_id`

* is required

* Type: `string`

* cannot be null

* defined in: [get-support-session output](get-support-session-output-properties-session_id.md "http://schema.nethserver.org/node/get-support-session-output.json#/properties/session_id")

### session\_id Type

`string`

## active

If the opened session is still active or not

`active`

* is required

* Type: `boolean`

* cannot be null

* defined in: [get-support-session output](get-support-session-output-properties-active.md "http://schema.nethserver.org/node/get-support-session-output.json#/properties/active")

### active Type

`boolean`
