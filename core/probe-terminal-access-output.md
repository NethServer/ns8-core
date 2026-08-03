# probe-terminal-access output Schema

```txt
http://schema.nethserver.org/node/probe-terminal-access-output.json
```

Whether sshd on this node accepts the cluster-admin terminal. Deliberately limited to three booleans and the port: the task output is readable by any authenticated user, so the raw sshd configuration must not be exposed.

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                         |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Forbidden             | none                | [probe-terminal-access-output.json](node/probe-terminal-access-output.json "open original schema") |

## probe-terminal-access output Type

`object` ([probe-terminal-access output](probe-terminal-access-output.md))

## probe-terminal-access output Examples

```json
{
  "permit_root_login": true,
  "password_auth": true,
  "listen_wg0": true,
  "port": 22
}
```

# probe-terminal-access output Properties

| Property                                  | Type      | Required | Nullable       | Defined by                                                                                                                                                                                       |
| :---------------------------------------- | :-------- | :------- | :------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [permit\_root\_login](#permit_root_login) | `boolean` | Required | cannot be null | [probe-terminal-access output](probe-terminal-access-output-properties-permit_root_login.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/permit_root_login") |
| [password\_auth](#password_auth)          | `boolean` | Required | cannot be null | [probe-terminal-access output](probe-terminal-access-output-properties-password_auth.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/password_auth")         |
| [listen\_wg0](#listen_wg0)                | `boolean` | Required | cannot be null | [probe-terminal-access output](probe-terminal-access-output-properties-listen_wg0.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/listen_wg0")               |
| [port](#port)                             | `integer` | Required | cannot be null | [probe-terminal-access output](probe-terminal-access-output-properties-port.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/port")                           |

## permit\_root\_login

True when root may authenticate with a password from the cluster VPN, i.e. the effective PermitRootLogin is yes

`permit_root_login`

* is required

* Type: `boolean`

* cannot be null

* defined in: [probe-terminal-access output](probe-terminal-access-output-properties-permit_root_login.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/permit_root_login")

### permit\_root\_login Type

`boolean`

## password\_auth

True when password authentication is enabled from the cluster VPN

`password_auth`

* is required

* Type: `boolean`

* cannot be null

* defined in: [probe-terminal-access output](probe-terminal-access-output-properties-password_auth.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/password_auth")

### password\_auth Type

`boolean`

## listen\_wg0

True when sshd listens on the node VPN address, or on every address

`listen_wg0`

* is required

* Type: `boolean`

* cannot be null

* defined in: [probe-terminal-access output](probe-terminal-access-output-properties-listen_wg0.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/listen_wg0")

### listen\_wg0 Type

`boolean`

## port

Effective sshd listening port

`port`

* is required

* Type: `integer`

* cannot be null

* defined in: [probe-terminal-access output](probe-terminal-access-output-properties-port.md "http://schema.nethserver.org/node/probe-terminal-access-output.json#/properties/port")

### port Type

`integer`

### port Constraints

**maximum**: the value of this number must smaller than or equal to: `65535`

**minimum**: the value of this number must greater than or equal to: `1`
