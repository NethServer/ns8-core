# set-password-warning Schema

```txt
http://schema.nethserver.org/cluster/set-password-warning-input.json
```

Input schema of the set-password-warning action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [set-password-warning-input.json](cluster/set-password-warning-input.json "open original schema") |

## set-password-warning Type

`object` ([set-password-warning](set-password-warning-input.md))

## set-password-warning Examples

```json
{
  "domain": "mydom.test",
  "notification": true,
  "days": 7,
  "mail_from": "no-reply@nethserver.org",
  "mail_template_name": "default_en"
}
```

```json
{
  "domain": "mydom.test",
  "notification": true,
  "days": 7,
  "mail_from": "no-reply@nethserver.org",
  "mail_template_name": "custom",
  "mail_template_content": "RGVhciAkbmFtZSAoJH...",
  "mail_subject": "Password expiration warning"
}
```

# set-password-warning Properties

| Property                                          | Type      | Required | Nullable       | Defined by                                                                                                                                                                                      |
| :------------------------------------------------ | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [domain](#domain)                                 | `string`  | Required | cannot be null | [set-password-warning](set-password-warning-input-properties-domain.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/domain")                               |
| [notification](#notification)                     | `boolean` | Required | cannot be null | [set-password-warning](set-password-warning-input-properties-notification.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/notification")                   |
| [mail\_template](#mail_template)                  | `string`  | Optional | cannot be null | [set-password-warning](set-password-warning-input-properties-mail_template.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_template")                 |
| [days](#days)                                     | `integer` | Optional | cannot be null | [set-password-warning](set-password-warning-input-properties-days.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/days")                                   |
| [mail\_from](#mail_from)                          | `string`  | Optional | cannot be null | [set-password-warning](set-password-warning-input-properties-mail_from.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_from")                         |
| [mail\_template\_content](#mail_template_content) | `string`  | Optional | cannot be null | [set-password-warning](set-password-warning-input-properties-mail_template_content.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_template_content") |
| [mail\_subject](#mail_subject)                    | `string`  | Optional | cannot be null | [set-password-warning](set-password-warning-input-properties-mail_subject.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_subject")                   |

## domain

A user domain name

`domain`

* is required

* Type: `string`

* cannot be null

* defined in: [set-password-warning](set-password-warning-input-properties-domain.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/domain")

### domain Type

`string`

### domain Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## notification

Enable or disable user domain notifications

`notification`

* is required

* Type: `boolean`

* cannot be null

* defined in: [set-password-warning](set-password-warning-input-properties-notification.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/notification")

### notification Type

`boolean`

## mail\_template

The notification template to use

`mail_template`

* is optional

* Type: `string`

* cannot be null

* defined in: [set-password-warning](set-password-warning-input-properties-mail_template.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_template")

### mail\_template Type

`string`

### mail\_template Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## days

The number of days before the expiration date to send the notification

`days`

* is optional

* Type: `integer`

* cannot be null

* defined in: [set-password-warning](set-password-warning-input-properties-days.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/days")

### days Type

`integer`

### days Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## mail\_from

The sender email address

`mail_from`

* is optional

* Type: `string`

* cannot be null

* defined in: [set-password-warning](set-password-warning-input-properties-mail_from.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_from")

### mail\_from Type

`string`

### mail\_from Constraints

**email**: the string must be an email address, according to [RFC 5322, section 3.4.1](https://tools.ietf.org/html/rfc5322 "check the specification")

## mail\_template\_content

The template content in base64 format, it can be a text file or an HTML file. It must respect python string.Template syntax. Required if mail\_template\_name is 'custom'

`mail_template_content`

* is optional

* Type: `string`

* cannot be null

* defined in: [set-password-warning](set-password-warning-input-properties-mail_template_content.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_template_content")

### mail\_template\_content Type

`string`

### mail\_template\_content Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## mail\_subject

The email subject. It must respect python string.Template syntax. Required if mail\_template\_name is 'custom'

`mail_subject`

* is optional

* Type: `string`

* cannot be null

* defined in: [set-password-warning](set-password-warning-input-properties-mail_subject.md "http://schema.nethserver.org/cluster/set-password-warning-input.json#/properties/mail_subject")

### mail\_subject Type

`string`

### mail\_subject Constraints

**minimum length**: the minimum number of characters for this string is: `1`
