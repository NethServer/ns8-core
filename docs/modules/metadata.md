---
layout: default
title: Metadata
nav_order: 32
parent: Modules
---

# Metadata

Local metadata are used to display information inside the *About* page.
Metadata are composed by:

- a JSON file named `metadata.json`, saved inside `imageroot/ui/public/metadata.json`
- a PNG image for the logo named `module_default_logo.png`, saved inside `imageroot/ui/src/assets/module_default_logo.png`

## metadata.json

The JSON format is very similar to the 
[repository metadata](https://github.com/NethServer/ns8-core/blob/main/core/imageroot/var/lib/nethserver/cluster/repodata-schema.json),
without `id`, `logo` and `versions` fields.

Example of `metadata.json`:
```json
{
  "name": "Kickstart",
  "upstream_name": "Kickstart 12",
  "description": {
    "en": "My kickstart module"
  },
  "categories": ["somecategory"],
  "authors": [
    {
      "name": "Name Surname",
      "email": "author@yourmail.org"
    }
  ],
  "docs": {
    "terms_url": "https://docs.kickstart.com/terms/",
    "documentation_url": "https://docs.kickstart.com/",
    "bug_url": "https://github.com/NethServer/dev",
    "code_url": "https://github.com/author/ns8-kickstart",
    "relnotes_url": "https://github.com/author/ns8-kickstart/releases"
  },
  "source": "ghcr.io/nethserver/kickstart"
}
```

The `categories` field contains a list of string. Each string is a category that should describe the application usage scenario.
It could be something like `database` or `collaboration`. Always use lower case strings without spaces.

The list of available categories is available in this [translation.json
source
file](https://github.com/NethServer/ns8-core/blob/main/core/ui/public/i18n/en/translation.json),
under the `app_categories` object. Change proposals to the category list can be
discussed on the community forum, as explained in [development process](../../development_process/).

The `relnotes_url` is an optional key in the module metadata that specifies the URL to the release changelog page. This URL points to the location where users can find detailed information about the changes, improvements, and fixes included in each release of the module.

## Logo

The logo must be a PNG image and its size should be at least of 128x128 pixels.
