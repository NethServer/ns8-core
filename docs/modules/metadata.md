---
layout: default
title: Metadata
nav_order: 15
parent: Modules
---

# Metadata

Local metadata are used to display information on the *About* page and the *Software Center* page.

Metadata consist of:

- JSON file: `metadata.json`
- Logo image: `module_default_logo.png`

## metadata.json

The file path, relative to the repository root, must be `imageroot/ui/public/metadata.json`.

The JSON format of this file is very similar to the [repository metadata](https://github.com/NethServer/ns8-core/blob/main/core/imageroot/var/lib/nethserver/cluster/repodata-schema.json), excluding the `id`, `logo`, and `versions` fields.

Example of `metadata.json`:

```json
{
  "name": "Kickstart App",
  "description": {
    "en": "My first Kickstart application"
  },
  "categories": [],
  "authors": [
    {
      "name": "Name Surname",
      "email": "author@example.org"
    }
  ],
  "docs": {
    "documentation_url": "http://docs.kickstart.example.org/",
    "bug_url": "https://bugs.kickstart.example.org",
    "code_url": "https://github.com/author/ns8-kickstart"
  },
  "source": "ghcr.io/author/kickstart"
}
```

* The `name` field is the human-readable name of the application. It is the string displayed to the UI user.
* The `description` object contains one attribute for each language. Values are localized strings in different languages. Attribute names must match the current language code used by the main Core app. See the function [getAppDescription()](https://github.com/search?) in the `ns8-ui-lib` repository for more details.
* The `source` field is the registry address of the application image.
* The `code_url` is the web location where the app's source code is published.
* The `bug_url` and `documentation_url` are the web locations for the application's public bug tracker and public documentation, respectively.
* The `categories` field contains a list of strings. Each string is a category that describes the application's usage scenario. Examples include `database` or `collaboration`. Always use lowercase strings without spaces. Allowed values are defined in the core repository, in the file [translation.json](https://github.com/NethServer/ns8-core/blob/2.8.2/core/ui/public/i18n/en/translation.json#L886).

## Logo image module_default_logo.png

The application logo must be a PNG image, saved inside `imageroot/ui/src/assets/module_default_logo.png`. It must be square-shaped and at least 128x128 pixels in size.
