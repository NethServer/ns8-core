---
layout: default
title: Metadata
nav_order: 15
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
  "name": "kickstart",
  "description": {
    "en": "My kickstart module"
  },
  "categories": [],
  "authors": [
    {
      "name": "Name Surname",
      "email": "author@yourmail.org"
    }
  ],
  "docs": {
    "documentation_url": "https://docs.kickstart.com/",
    "bug_url": "https://github.com/NethServer/dev",
    "code_url": "https://github.com/author/ns8-kickstart"
  },
  "source": "ghcr.io/nethserver/kickstart"
}
```

The `categories` field contains a list of string. Each string is a category that should describe the application usage scenario.
It could be something like `database` or `collaboration`. Always use lower case strings without spaces.

## Logo

The logo must be a PNG image and its size should be at least of 128x128 pixels.
