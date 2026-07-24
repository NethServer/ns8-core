# Output for list-installed-modules Schema

```txt
http://schema.nethserver.org/cluster/list-installed-modules-output.json
```

list-installed-modules output

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                                              |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :------------------------------------------------------------------------------------------------------ |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-installed-modules-output.json](cluster/list-installed-modules-output.json "open original schema") |

## Output for list-installed-modules Type

`object` ([Output for list-installed-modules](list-installed-modules-output.md))

## Output for list-installed-modules Examples

```json
{
  "ghcr.io/nethserver/traefik": [
    {
      "id": "traefik1",
      "node": "1",
      "digest": "sha256:64a4ccd3c5ded935ac8d28d2bd534e1b305b80a131a3fe3fe5e123c03b4aa772",
      "source": "ghcr.io/nethserver/traefik",
      "version": "latest",
      "logo": "https://raw.githubusercontent.com/NethServer/ns8-repomd/repomd/dokuwiki/logo.png",
      "module": "traefik",
      "flags": [
        "no_data_backup",
        "core_module"
      ]
    }
  ],
  "ghcr.io/nethserver/dokuwiki": [
    {
      "id": "dokuwiki2",
      "node": "1",
      "digest": "sha256:929465c177d9e40559b3fa838f2e429060ff7e8c3e4e1a12076b3304ad562982",
      "source": "ghcr.io/nethserver/dokuwiki",
      "version": "0.0.1-alpha1",
      "logo": "https://raw.githubusercontent.com/NethServer/ns8-repomd/repomd/dokuwiki/logo.png",
      "module": "dokuwiki",
      "flags": []
    },
    {
      "id": "dokuwiki1",
      "node": "1",
      "digest": "sha256:929465c177d9e40559b3fa838f2e429060ff7e8c3e4e1a12076b3304ad562982",
      "source": "ghcr.io/nethserver/dokuwiki",
      "version": "0.0.1-alpha1",
      "logo": "https://raw.githubusercontent.com/NethServer/ns8-repomd/repomd/dokuwiki/logo.png",
      "module": "dokuwiki",
      "flags": []
    }
  ]
}
```

# Output for list-installed-modules Properties

| Property | Type    | Required | Nullable       | Defined by                                                                                                                                                                               |
| :------- | :------ | :------- | :------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `.*`     | `array` | Optional | cannot be null | [Output for list-installed-modules](list-installed-modules-output-patternproperties-.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*") |

## Pattern: `.*`



`.*`

* is optional

* Type: `object[]` ([Details](list-installed-modules-output-patternproperties--items.md))

* cannot be null

* defined in: [Output for list-installed-modules](list-installed-modules-output-patternproperties-.md "http://schema.nethserver.org/cluster/list-installed-modules-output.json#/patternProperties/.*")

### .\* Type

`object[]` ([Details](list-installed-modules-output-patternproperties--items.md))
