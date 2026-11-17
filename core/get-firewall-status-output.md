# get-firewall-status output Schema

```txt
http://schema.nethserver.org/node/get-firewall-status-output.json
```

Output schema of the get-firewall-status action

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                     |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :--------------------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [get-firewall-status-output.json](node/get-firewall-status-output.json "open original schema") |

## get-firewall-status output Type

unknown ([get-firewall-status output](get-firewall-status-output.md))

## get-firewall-status output Examples

```json
{
  "interfaces": [
    "eth0",
    "eth1"
  ],
  "services": [
    {
      "name": "dhcpv6-client",
      "ports": [
        "546/udp"
      ]
    },
    {
      "name": "https",
      "ports": [
        "443/tcp"
      ]
    }
  ],
  "ports": [
    "55820/udp"
  ]
}
```
