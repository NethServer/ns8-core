# Untitled string in alter-backup-repository input Schema

```txt
http://schema.nethserver.org/cluster/alter-backup-repository-input.json#/$defs/smb_parameters/properties/smb_host
```

Host name or IP address

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                                |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [alter-backup-repository-input.json\*](cluster/alter-backup-repository-input.json "open original schema") |

## smb\_host Type

`string`

## smb\_host Constraints

**minimum length**: the minimum number of characters for this string is: `1`

**hostname**: the string must be a hostname, according to [RFC 1123, section 2.1](https://tools.ietf.org/html/rfc1123 "check the specification")
