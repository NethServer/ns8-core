# get-fqdn output Schema

```txt
http://schema.nethserver.org/node/get-fqdn-output.json
```

Output schema of the get-fqdn action

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                               |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :----------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [get-fqdn-output.json](node/get-fqdn-output.json "open original schema") |

## get-fqdn output Type

`object` ([get-fqdn output](get-fqdn-output.md))

## get-fqdn output Examples

```json
{
  "hostname": "dn2",
  "domain": "example.org"
}
```

# get-fqdn output Properties

| Property              | Type     | Required | Nullable       | Defined by                                                                                                                              |
| :-------------------- | :------- | :------- | :------------- | :-------------------------------------------------------------------------------------------------------------------------------------- |
| [hostname](#hostname) | `string` | Required | cannot be null | [get-fqdn output](get-fqdn-output-properties-hostname.md "http://schema.nethserver.org/node/get-fqdn-output.json#/properties/hostname") |
| [domain](#domain)     | `string` | Required | cannot be null | [get-fqdn output](get-fqdn-output-properties-domain.md "http://schema.nethserver.org/node/get-fqdn-output.json#/properties/domain")     |

## hostname



`hostname`

* is required

* Type: `string`

* cannot be null

* defined in: [get-fqdn output](get-fqdn-output-properties-hostname.md "http://schema.nethserver.org/node/get-fqdn-output.json#/properties/hostname")

### hostname Type

`string`

## domain



`domain`

* is required

* Type: `string`

* cannot be null

* defined in: [get-fqdn output](get-fqdn-output-properties-domain.md "http://schema.nethserver.org/node/get-fqdn-output.json#/properties/domain")

### domain Type

`string`
