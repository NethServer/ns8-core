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

It's important to create a release on GitHub for each new version of the module,
because the release allows to attach the SBOM files and create a changelog.

The release process is divided into two main steps:

1. tag the repository with the new version, following the [version numbering](https://handbook.nethserver.org/version_numbering/) guidelines
2. create a release on GitHub, attaching the SBOM files generated during the build process

The best and recommended way to create a new release is using the [GitHub CLI](https://cli.github.com/).
After installing the tool, use the following command:

```bash
gh release create
```

The command will create a tag and a release on GitHub.
In case of stable releases, the command will also attach the SBOM files generated during the build process.

Then, you can sync the remote tag with the local repository using the following command:

```bash
git fetch --tags
```

For testing releases, you can also use [ns8-release-module](https://github.com/NethServer/gh-ns8-release-module/). 
It is a [GitHub CLI](https://cli.github.com/) Extension that automates the creation and management and it's quite useful when
creating testing releases.
See the full documentation inside the [README](https://github.com/NethServer/gh-ns8-release-module/blob/main/Readme.md).

### Manual release

**Warning**: always prefer the GitHub CLI to create a release.

In case you forgot to create a release using the GitHub CLI, you can create it manually.
First, make sure the tag is created in the repository and pushed to the remote. Then:

1. Navigate to the GitHub UI:
    - Go to the repository's "Releases" section.
    - Click on "Draft a new release" to start the process.

3. Add a changelog:
    - Provide a detailed changelog for the release.
    - Optionally, use the automatic changelog generation button if available.

4. Attach the SBOM file:
    - For stable releases, attach the Software Bill of Materials (SBOM) file to the release by respecting conventions describe in the [SBOM generation](https://handbook.nethserver.org/security/#sbom-software-bill-of-materials) page.
    You can find the SBOM files as artifacts of the GitHub Actions workflow.