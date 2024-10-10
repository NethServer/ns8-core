---
layout: default
title: Automatic builds
nav_order: 132
parent: Modules
---

## Automatic builds

Every time a commit is pushed to the repository, a new build can
be automatically started.
To enable automatic builds, create a YAML file like `.github/workflows/<module>.yaml`.

Use [dokuwki Github workflow](https://github.com/NethServer/ns8-dokuwiki/tree/main/.github/workflows) as a template, make sure to replace all occurrences of `dokuwiki`
with the name of the new module.

Upon merging into the master branch, also the following automatic actions will be started:
- metadata build on [ns8-repomd](https://github.com/NethServer/ns8-repomd/)
- api documentation build on [apidoc branch](https://github.com/NethServer/ns8-core/tree/apidoc) (see also [API](/api)) 
