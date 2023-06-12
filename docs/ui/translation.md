---
layout: default
title: Localization
nav_order: 5
parent: User Interface
---

## Localization

The user interface is translated with [Weblate](https://hosted.weblate.org/projects/ns8/).

Every time a ns8 repository is updated, Github will notify Weblate.
When new translations are avaialable, Weblate will create a new pull request on the target repository.

To setup the translation process:

1. add [GitHub Weblate app](https://docs.weblate.org/en/latest/admin/continuous.html#github-setup) to your repository
2. add your repository to [hosted.weblate.org](https://hosted.weblate.org) or ask a NethServer developer to add it to ns8 Weblate project

You can automate the step 2 by using the following `weblate-add` script.

```bash
set -e
name=$1

curl \
    --fail \
    --data-binary '{
        "branch": "main",
        "file_format": "json-nested",
        "filemask": "ui/public/i18n/*/translation.json",
        "name": "'$name'",
        "slug": "'$name'",
        "repo": "https://github.com/NethServer/ns8-'$name'.git",
        "template": "ui/public/i18n/en/translation.json",
        "new_base": "",
        "vcs": "github",
        "license": "GPL-3.0-or-later"
    }' \
    -H "Content-Type: application/json" \
    -H "Authorization: Token $TOKEN" \
    https://hosted.weblate.org/api/projects/ns8/components/

curl \
    --fail \
    -d '{"language_code": "it"}' \
    -H "Content-Type: application/json" \
    -H "Authorization: Token $TOKEN" \
    https://hosted.weblate.org/api/components/ns8/$name/translations/
```

Usage example, where `TOKEN` is your Weblate token:
```
TOKEN=xxxxxxxxxxxx bash weblate-add nextcloud
```
