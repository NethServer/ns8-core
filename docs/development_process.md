---
layout: default
title: Development process
nav_order: 3
---

# Development process

All NethServer projects follow a shared development process.

See [NethServer development handbook](https://handbook.nethserver.org/) for a detailed guide on how to contribute to NethServer.

Here you can find some additional guidelines regarding the releases of modules, including the core.

## Release process

The release process is documented in the [release process](https://handbook.nethserver.org/release-process/) page.
But some additional steps are required for a correct release of all NS8 modules, including the core.

The release process is divided into two main steps:

1. tag the repository with the new version, following the [version numbering](https://handbook.nethserver.org/version_numbering/) guidelines
2. create a release on GitHub, attaching the SBOM files generated during the build process

There are two different ways to create a release:

1. manual mode: the developer creates the release on GitHub and eventually attaches the SBOM files
2. automatic mode: the release is created using a GitHub CLI extension

### Manual release

First, make sure the tag is created in the repository and pushed to the remote. Then:

1. Navigate to the GitHub UI:
    - Go to the repository's "Releases" section.
    - Click on "Draft a new release" to start the process.

3. Add a changelog:
    - Provide a detailed changelog for the release.
    - Optionally, use the automatic changelog generation button if available.

4. Attach the SBOM file:
    - For stable releases, attach the Software Bill of Materials (SBOM) file to the release by respecting conventions describe in the [SBOM generation](https://handbook.nethserver.org/security/#sbom-software-bill-of-materials) page.

### Automatic release (recommended)

The [ns8-release-module](https://github.com/NethServer/gh-ns8-release-module/) is a [GitHub CLI](https://cli.github.com/) Extension that automates the creation and management of
NethServer 8 modules releases.

See how to install and use the extension inside the [README](https://github.com/NethServer/gh-ns8-release-module/blob/main/Readme.md).