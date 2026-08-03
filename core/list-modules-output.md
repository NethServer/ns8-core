# list-modules output Schema

```txt
http://schema.nethserver.org/cluster/list-modules-output.json
```

List modules output

| Abstract            | Extensible | Status         | Identifiable            | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                          |
| :------------------ | :--------- | :------------- | :---------------------- | :---------------- | :-------------------- | :------------------ | :---------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | Unknown identifiability | Forbidden         | Allowed               | none                | [list-modules-output.json](cluster/list-modules-output.json "open original schema") |

## list-modules output Type

`object[]` ([Details](list-modules-output-items.md))

## list-modules output Examples

```json
[
  {
    "name": "dokuwiki",
    "description": {
      "en": "Auto-generated description for traefik"
    },
    "logo": "http://127.0.0.1:8000/logo.png",
    "screenshots": [
      "http://127.0.0.1:8000/screenshots/template%3Adokuwiki_template.png"
    ],
    "categories": [
      "unknown"
    ],
    "authors": [
      {
        "name": "unknown",
        "email": "info@nethserver.org"
      }
    ],
    "docs": {
      "documentation_url": "https://docs.nethserver.org",
      "bug_url": "https://github.com/NethServer/dev",
      "code_url": "https://github.com/NethServer/"
    },
    "source": "ghcr.io/nethserver/dokuwiki",
    "versions": [
      {
        "tag": "0.0.2",
        "testing": false,
        "labels": {
          "io.buildah.version": "1.19.6",
          "org.nethserver.rootfull": "0",
          "org.nethserver.tcp_ports_demand": "1"
        }
      },
      {
        "tag": "0.0.1",
        "testing": false,
        "labels": {
          "io.buildah.version": "1.19.6",
          "org.nethserver/rootfull": "0",
          "org.nethserver/tcp_ports_demand": "1"
        }
      },
      {
        "tag": "0.0.1-alpha1",
        "testing": true,
        "labels": {
          "io.buildah.version": "1.19.6",
          "org.nethserver/rootfull": "0",
          "org.nethserver/tcp_ports_demand": "1"
        }
      }
    ],
    "repository": "t3",
    "repository_updated": "Mon, 28 Jun 2021 14:42:44 GMT",
    "install_destinations": [
      {
        "node_id": 1,
        "instances": 2,
        "eligible": true,
        "reject_reason": null
      }
    ],
    "certification_level": 2,
    "rootfull": false,
    "updates": [
      {
        "id": "dokuwiki2",
        "node": "1",
        "digest": "sha256:929465c177d9e40559b3fa838f2e429060ff7e8c3e4e1a12076b3304ad562982",
        "source": "ghcr.io/nethserver/dokuwiki",
        "version": "0.0.1-alpha1",
        "update": "0.0.2"
      }
    ],
    "disabled_updates_reason": "",
    "installed": [
      [
        {
          "id": "dokuwiki2",
          "node": "1",
          "digest": "sha256:929465c177d9e40559b3fa838f2e429060ff7e8c3e4e1a12076b3304ad562982",
          "source": "ghcr.io/nethserver/dokuwiki",
          "version": "0.0.1-alpha1"
        },
        {
          "id": "dokuwiki1",
          "node": "1",
          "digest": "sha256:929465c177d9e40559b3fa838f2e429060ff7e8c3e4e1a12076b3304ad562982",
          "source": "ghcr.io/nethserver/dokuwiki",
          "version": "0.0.1-alpha1"
        }
      ]
    ]
  }
]
```
