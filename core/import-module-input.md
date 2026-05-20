# import-module input Schema

```txt
http://schema.nethserver.org/module/import-module-input.json
```

Import the module state from an external procedure with rsync

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                         |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [import-module-input.json](module/import-module-input.json "open original schema") |

## import-module input Type

`object` ([import-module input](import-module-input.md))

## import-module input Examples

```json
{
  "credentials": [
    "dokuwiki1",
    "s3cr3t"
  ],
  "port": 20027,
  "volumes": [
    "dokuwiki-data"
  ]
}
```

```json
{
  "credentials": [
    "samba3",
    "s3cr3t"
  ],
  "port": 20028,
  "volumes": [
    "data",
    "config"
  ]
}
```

# import-module input Properties

| Property                    | Type      | Required | Nullable       | Defined by                                                                                                                                                                                      |
| :-------------------------- | :-------- | :------- | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [volumes](#volumes)         | `array`   | Required | cannot be null | [import-module input](import-module-input-properties-initial-volume-set-where-the-module-state-is-stored.md "http://schema.nethserver.org/module/import-module-input.json#/properties/volumes") |
| [port](#port)               | `integer` | Required | cannot be null | [import-module input](import-module-input-properties-rsyncd-tcp-port-number.md "http://schema.nethserver.org/module/import-module-input.json#/properties/port")                                 |
| [credentials](#credentials) | `array`   | Required | cannot be null | [import-module input](import-module-input-properties-rsyncd-service-credentials.md "http://schema.nethserver.org/module/import-module-input.json#/properties/credentials")                      |

## volumes



`volumes`

* is required

* Type: `string[]`

* cannot be null

* defined in: [import-module input](import-module-input-properties-initial-volume-set-where-the-module-state-is-stored.md "http://schema.nethserver.org/module/import-module-input.json#/properties/volumes")

### volumes Type

`string[]`

## port



`port`

* is required

* Type: `integer` ([Rsyncd TCP port number](import-module-input-properties-rsyncd-tcp-port-number.md))

* cannot be null

* defined in: [import-module input](import-module-input-properties-rsyncd-tcp-port-number.md "http://schema.nethserver.org/module/import-module-input.json#/properties/port")

### port Type

`integer` ([Rsyncd TCP port number](import-module-input-properties-rsyncd-tcp-port-number.md))

### port Constraints

**minimum**: the value of this number must greater than or equal to: `1`

## credentials

An array with two elements: username, password

`credentials`

* is required

* Type: `string[]`

* cannot be null

* defined in: [import-module input](import-module-input-properties-rsyncd-service-credentials.md "http://schema.nethserver.org/module/import-module-input.json#/properties/credentials")

### credentials Type

`string[]`

### credentials Constraints

**maximum number of items**: the maximum number of items for this array is: `2`

**minimum number of items**: the minimum number of items for this array is: `2`
