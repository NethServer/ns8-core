# get smarthost settings Schema

```txt
http://schema.nethserver.org/cluster/get-smarthost.json
```

Get the settings an external smarthost provider

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                              |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-smarthost.json](cluster/get-smarthost.json "open original schema") |

## get smarthost settings Type

`object` ([get smarthost settings](get-smarthost.md))

## get smarthost settings Examples

```json
{
  "host": "smtp.example.com",
  "port": 587,
  "username": "username",
  "password": "secret",
  "enabled": true,
  "encrypt_smtp": "starttls",
  "tls_verify": false,
  "mail_server": [
    {
      "mail_id": "mail6",
      "mail_name": "toto.com",
      "node": "1",
      "host": "10.5.4.1",
      "node_name": "node1"
    }
  ],
  "manual_configuration": false
}
```

# get smarthost settings Properties

| Property                                       | Type      | Required | Nullable       | Defined by                                                                                                                                                            |
| :--------------------------------------------- | :-------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [username](#username)                          | `string`  | Required | cannot be null | [get smarthost settings](get-smarthost-properties-username.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/username")                         |
| [port](#port)                                  | `integer` | Required | cannot be null | [get smarthost settings](get-smarthost-properties-port.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/port")                                 |
| [password](#password)                          | `string`  | Required | cannot be null | [get smarthost settings](get-smarthost-properties-password.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/password")                         |
| [host](#host)                                  | `string`  | Required | cannot be null | [get smarthost settings](get-smarthost-properties-host.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/host")                                 |
| [enabled](#enabled)                            | `boolean` | Required | cannot be null | [get smarthost settings](get-smarthost-properties-enabled.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/enabled")                           |
| [encrypt\_smtp](#encrypt_smtp)                 | `string`  | Required | cannot be null | [get smarthost settings](get-smarthost-properties-encrypt_smtp.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/encrypt_smtp")                 |
| [tls\_verify](#tls_verify)                     | `boolean` | Required | cannot be null | [get smarthost settings](get-smarthost-properties-tls_verify.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/tls_verify")                     |
| [mail\_server](#mail_server)                   | `array`   | Required | cannot be null | [get smarthost settings](get-smarthost-properties-mail_server.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/mail_server")                   |
| [manual\_configuration](#manual_configuration) | `boolean` | Required | cannot be null | [get smarthost settings](get-smarthost-properties-manual_configuration.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/manual_configuration") |

## username

Username to use the smarthost smtp

`username`

* is required

* Type: `string` ([username](get-smarthost-properties-username.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-username.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/username")

### username Type

`string` ([username](get-smarthost-properties-username.md))

## port

Port number of the smtp smarthost

`port`

* is required

* Type: `integer` ([port](get-smarthost-properties-port.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-port.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/port")

### port Type

`integer` ([port](get-smarthost-properties-port.md))

### port Constraints

**maximum**: the value of this number must smaller than or equal to: `65535`

**minimum**: the value of this number must greater than or equal to: `1`

## password

Password to use the smarthost

`password`

* is required

* Type: `string` ([password](get-smarthost-properties-password.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-password.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/password")

### password Type

`string` ([password](get-smarthost-properties-password.md))

## host

Host name for the smarthost, like 'wiki.nethserver.org'

`host`

* is required

* Type: `string` ([host](get-smarthost-properties-host.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-host.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/host")

### host Type

`string` ([host](get-smarthost-properties-host.md))

### host Constraints

**(international) hostname**: the string must be an (IDN) hostname, according to [RFC 5890, section 2.3.2.3](https://tools.ietf.org/html/rfc5890 "check the specification")

## enabled

Enable or disable the smarthost settings

`enabled`

* is required

* Type: `boolean` ([enabled](get-smarthost-properties-enabled.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-enabled.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/enabled")

### enabled Type

`boolean` ([enabled](get-smarthost-properties-enabled.md))

## encrypt\_smtp

Enable or disable the tls encryption with the smtp server

`encrypt_smtp`

* is required

* Type: `string` ([encrypt\_smtp](get-smarthost-properties-encrypt_smtp.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-encrypt_smtp.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/encrypt_smtp")

### encrypt\_smtp Type

`string` ([encrypt\_smtp](get-smarthost-properties-encrypt_smtp.md))

### encrypt\_smtp Constraints

**enum**: the value of this property must be equal to one of the following values:

| Value        | Explanation |
| :----------- | :---------- |
| `"none"`     |             |
| `"starttls"` |             |
| `"tls"`      |             |

## tls\_verify

Enable or disable the usage of a valid certificate

`tls_verify`

* is required

* Type: `boolean` ([tls\_verify](get-smarthost-properties-tls_verify.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-tls_verify.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/tls_verify")

### tls\_verify Type

`boolean` ([tls\_verify](get-smarthost-properties-tls_verify.md))

## mail\_server



`mail_server`

* is required

* Type: `object[]` ([Details](get-smarthost-properties-mail_server-items.md))

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-mail_server.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/mail_server")

### mail\_server Type

`object[]` ([Details](get-smarthost-properties-mail_server-items.md))

## manual\_configuration



`manual_configuration`

* is required

* Type: `boolean`

* cannot be null

* defined in: [get smarthost settings](get-smarthost-properties-manual_configuration.md "http://schema.nethserver.org/cluster/get-smarthost.json#/properties/manual_configuration")

### manual\_configuration Type

`boolean`
