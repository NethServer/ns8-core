---
layout: default
title: Build system and SBOM
nav_order: 10
---

# Build system

* TOC
{:toc}

## Manual builds

Every module, including the core, can be built locally or automatically using GitHub Actions.
Each repository usually has a file named `build-images.sh` that can be used to build the module images.

When executed on a local machines, it requires the `podman` or `buildah` command line tools installed.
To build the images, execute the `build-images.sh` script.

Once the images are built, they are available on the local machine and can be used for testing.
To list the images, use the `podman images` command.
The images can be pushed to the image registry using the `podman push` command.

The same script is used by the GitHub Actions workflow to build the images and push them to the image registry.
When executed inside the GitHub Actions environment, the script outputs the image names that are
automatically pushed to `ghcr.io` image registry.

## Automatic builds

Every time a commit is pushed to the repository, a new build can
be automatically started.

Each module, including the core, should use the standard GitHub Actions workflow name `publish-branch.yml` available inside the
[ns8-github-actions](https://github.com/NethServer/ns8-github-actions/) repository.

To enable automatic builds, create a YAML file named `publish-images.yaml` inside `.github/workflows/` directory,
like `.github/workflows/publish-image.yaml`.

Copy this content into the file:

```yaml
name: "Publish images"

on:
  push:
  workflow_dispatch:

permissions:
  packages: write

jobs:
  publish-images:
    if: github.run_number > 1
    uses: NethServer/ns8-github-actions/.github/workflows/publish-branch.yml@v1
```

## SBOM generation

A SBOM (Software Bill of Materials) is a complete list of all software parts that make up a product.
See the developer handbook [security](https://handbook.nethserver.org/security/) for more information.

SBOM are part of the build process and must be generated when a new stable version of the core or a module is released.
It can be created using the `scan-with-trivy.yml` action available in the [ns8-github-actions](https://github.com/NethServer/ns8-github-actions/) repository.

The workflow provides flexibility with configurable inputs, enabling targeted
scanning of specific components (`ui`, `imageroot`) or container images. It
supports key actions such as vulnerability detection, updating the GitHub
Dependency Graph, and generating CycloneDX SBOMs, which are uploaded as
artifacts or attached to GitHub releases when applicable.

See the [release process](../development_process#release-process/) page for more information on how to create a release.

**Inputs:**

1. **`images`**: A JSON array representing the list of container
   images to scan.
2. **`scan_module_ui`**: A boolean flag to enable/disable scanning of the module UI
   (default: `true`).
3. **`scan_module_imageroot`**: A boolean flag to enable/disable scanning of
   the module imageroot (default: `true`).
4. **`generate_sbom`**: A boolean flag to enable/disable SBOM generation
   (default: `true`).
5. **`update_dependencies_graph`**: A boolean flag to enable/disable the
   updating of the GitHub Dependency Graph (default: `true`).
6. **`vulnerability_scan`**: A boolean flag to enable/disable vulnerability
   scanning (default: `true`).
7. **`severity`**: A comma-separated string of severity levels to include in
   the scan results (default: `"UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL"`).

**Outputs:**

1. **Vulnerability Scan Results**:
   - Format: SARIF
   - Location: Uploaded to the GitHub Security tab.

2. **SBOM Files**:
   - Format: CycloneDX (JSON), named with the format `module.cdx.json`.
   - Location:
     - Uploaded as artifacts to the workflow run.
     - Optionally attached to GitHub releases if a tag triggers the workflow.

3. **Dependency Graph Updates**:
   - Format: Github Dependency Graph
   - Dependency data is pushed to the GitHub Dependency Graph to provide a view
     of current dependencies.

### Usage

The workflow can be used in two ways for running Trivy security scans. It can
be simply included as a job into the existing "Publish images" workflow, where
it runs after other jobs finish up. Or it can be set up as a separate workflow
that automatically triggers when image publishing is complete. 

Both approaches are effective - the first keeps everything together in one
pipeline, while the second allows for independent management of security
scanning separate from the main build process. 

However, one disadvantage of the second method is that the workflow will always
appear as executed in the GitHub Actions list after the publish workflow, even
when the scan is not actually performed. This can create misleading entries in
the workflow history compared to the integrated approach.

In both methods the job runs if the output from the `module` job indicates that
the release is either `stable` or `latest`, or if the workflow was manually
triggered (`workflow_dispatch`). It passes two parameters to this workflow: the
images to be scanned, which are outputs from the `module` job, and the severity
levels to be reported.

#### Generate SBOM from build `Publish images` workflow (recommended)

Call the `scan-with-trivy.yml` action as a job inside the `Publish images` workflow.

If the `Publish images` workflow is already present in the repository, make sure to:

- update the permissions section to include the necessary write permissions
- add the `module` job to the workflow, data collected from this job is used in the `trivy` job
- add the `trivy` job to the workflow

{% raw %}
```yaml
name: "Publish images"

on:
  push:
  workflow_dispatch:

permissions:
  packages: write
  actions: read
  contents: write
  security-events: write

jobs:
  publish-images:
    if: github.run_number > 1
    uses: NethServer/ns8-github-actions/.github/workflows/publish-branch.yml@v1
  module:
    needs: publish-images
    uses: NethServer/ns8-github-actions/.github/workflows/module-info.yml@v1
  trivy:
    needs: module
    if: ${{ needs.module.outputs.release == 'stable' || needs.module.outputs.release == 'latest' || github.event_name == 'workflow_dispatch' }}
    uses: NethServer/ns8-github-actions/.github/workflows/scan-with-trivy.yml@v1
    with:
      images: ${{ needs.module.outputs.images }}
```
{% endraw %}

Remember to commit and push the changes to the repository.

#### Generate SBOM using a dedicated workflow

Inside the `.github/workflows` directory, create a new YAML file, `analyze-module.yml`.
Copy and paste the YAML configuration into the `analyze-module.yml` file:

{% raw %}
```yaml
name: Analyze module

on:
  workflow_dispatch:
  workflow_run:
    workflows: ["Publish images"]
    types:
      - completed

permissions:
  packages: write
  actions: read
  contents: write
  security-events: write

jobs:
  module:
    if: ${{ github.event.workflow_run.conclusion == 'success' || github.event.workflow_run.conclusion == '' }}
    uses: NethServer/ns8-github-actions/.github/workflows/module-info.yml@v1
  trivy:
    needs: module
    if: ${{ needs.module.outputs.release == 'stable' || needs.module.outputs.release == 'latest' || github.event_name == 'workflow_dispatch' }}
    uses: NethServer/ns8-github-actions/.github/workflows/scan-with-trivy.yml@v1
    with:
      images: ${{ needs.module.outputs.images }}
      severity: "HIGH,CRITICAL"
```
{% endraw %}

Remember to commit and push the changes to the repository.

