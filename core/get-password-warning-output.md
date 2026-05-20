# get-password-warning-output Schema

```txt
http://schema.nethserver.org/cluster/get-password-warning-output.json
```

Output schema of the get-password-warning action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [get-password-warning-output.json](cluster/get-password-warning-output.json "open original schema") |

## get-password-warning-output Type

unknown ([get-password-warning-output](get-password-warning-output.md))

## get-password-warning-output Examples

```json
{
  "domain": "mydom.test",
  "notification": true,
  "days": 7,
  "mail_from": "no-reply@nethserver.org",
  "mail_template_name": "default_en",
  "smtp_enabled": false
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
  "mail_subject": "Password expiration warning",
  "smtp_enabled": true
}
```
