---
layout: default
title: Development process
nav_order: 3
---

# Development process

The NethServer project follows a shared development process. See
[NethServer developer handbook](https://handbook.nethserver.org/) for a
detailed guide on **how to contribute to NethServer** and how we keep
track of the work with GitHub Project and Issues.

The following section describes the release process of NS8 core
and application images.

## NS8 release process

The release process is fully automated with a [GitHub Actions
workflow][^gha]. The workflow is triggered when a [Semver tag][^semver] is pushed
to the git repository.

[^gha]: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax
[^semver]: https://handbook.nethserver.org/version_numbering/#version-numbering-rules

An action runner starts and executes the main release procedure as defined
by a .yml file under ``.github/workflows/`` directory:

- `publish-images.yml`, for application images,
- `core.yml`, for core images.

In both cases the repository build script (e.g. `build-images.sh`) is
executed by the action runner and resulting Podman container images are
uploaded to the `ghcr.io` image registry.

There may be other conditional workflows that follow. For example:

- If the Semver tag is _stable_, the `ns8-core` repository also builds and
  uploads .qcow and .vmdk images to DigitalOcean Spaces.
- If the stable tag corresponds to a GitHub release, SBOM signatures are
  produced and attached to the release.

The repository GitHub release page is referenced by NS8 Software Center UI
to provide a human-readable summary of changes (changelog). Hence **a
GitHub release must be created**, with the same name as the Semver tag.

Since SBOM signatures and the release changelog are mandatory, the
following section explains how to simultaneously create the repository tag
with a GitHub release.

### Stable release with `gh release` command

As pre-requisite, install `gh`, the official [GitHub CLI][^ghcli] command,
to trigger the release process from a terminal.

[^ghcli]: https://cli.github.com/

#### 1. Last release number

Find the last release number, to use as base to inspect the pending
changes on the `main` branch and generate the changelog.

    ⬢ [davidep@toolbx ns8-core]$ gh release list -L 3
    TITLE   TYPE    TAG NAME  PUBLISHED
    3.17.1  Latest  3.17.1    about 6 days ago
    3.17.0          3.17.0    about 21 days ago
    3.16.1          3.16.1    about 24 days ago

For example, latest is `3.17.1`.

#### 2. Inspect pending changes

The goal of this step is to ensure that all pending changes have gone
through the QA review process.

Switch to `main` branch and pull latest changes.

    git checkout main
    git pull origin main

Review the pending changes changelog:

    git log
    git diff --stat 3.17.1
    git diff 3.17.1

Commits often contain a reference to pull requests. Reference to the
related issue is reported in the PR description.

#### 3. Find related issues

Find in the issue tracker the items with
[VERIFIED][^vrfy] label.
Filter by project, e.g. `NethServer`, `NethVoice` or by milestone.

[^vrfy]: https://github.com/NethServer/dev/issues?q=is%3Aissue%20state%3Aopen%20label%3Averified

If there are still non-verified commits, for example commits that belong
to issues in development or QA stage, **the release process must be stopped**.

#### 4. Check documentation

If a related issue requires a documentation update, make sure the
documentation changes are ready for publication in [ns8-docs
repository][^ns8docs].

[^ns8docs]: https://github.com/NethServer/ns8-docs

#### 5. Check translations

If there are pending UI changes, make sure related strings are ready for
translation on [NS8 Weblate project][^weblate].

[^weblate]: https://hosted.weblate.org/projects/ns8/#components

Check also if a pull request from Weblate is ready, and merge it.

#### 6. Create the release

Choose the [release version number][^semver], incrementing the latest
version major, minor, or patch number.  For example `3.17.1 -> 3.18.0` is
a minor increment.

Create the release and git tag on GitHub in a single command:

    gh release create 3.18.0 --generate-notes --notes-start-tag 3.17.1

Note that if you run `gh` on a git local repository, the working tree
state does not matter.  The command creates the tag on the last commit of
the main repository branch. Use ``--target`` to select an alternative
commit, e.g.:

    gh release create 3.18.0 --generate-notes --notes-start-tag 3.17.1 --target COMMIT-SHA

Then, you can sync the remote tag with the local repository using the
following command:

    git fetch --tags

#### 7. Publish documentation

To publish a documentation PR found at step 4, just merge it.

#### 8. Verification checks

1. The GitHub release has the list of closed PRs (changelog)
1. The SBOM artifacts were attached to the release
1. The container images were pushed to GitHub Packages (ghcr.io)
1. For ns8-core only, the VM pre-built images were uploaded to
   DigitalOcean Spaces
1. Documentation was published
1. Translation repository is not locked
1. Related issues were closed with a link to the GitHub release

## Testing releases

Testing releases occur when code changes are merged with the `main`
branch, to start a testing (QA) process.

Choose the [release version number][^semver], incrementing the latest
stable version major, minor, or patch number. For example, `3.17.1 ->
3.17.2` is a patch increment for a bug fix. The latest release is visible
in the GitHub repository page.

Then append a pre-release suffix, like `-testing.1`, or `-dev.1`:

    3.17.2-testing.1

Create the testing tag on the last commit of your working tree:

    git tag 3.17.2-testing.1

Push the tag to trigger the testing build:

    git push origin 3.17.2-testing.1

For testing releases, you can also use
[ns8-release-module](https://github.com/NethServer/gh-ns8-release-module/).

It is a [GitHub CLI](https://cli.github.com/) extension that automates the
creation and management and it's quite useful when creating testing
releases.

See the full documentation inside the
[README](https://github.com/NethServer/gh-ns8-release-module/blob/main/Readme.md).

## Create a release from GitHub

> **Warning**: always prefer the GitHub CLI to create a release.

In case you forgot to create a release using the GitHub CLI, you can create it manually.
First, make sure the tag is created in the repository and pushed to the remote. Then:

1. Navigate to the GitHub UI:
    - Go to the repository's "Releases" section.
    - Click on "Draft a new release" to start the process.

2. Add a changelog:
    - Provide a detailed changelog for the release.
    - Optionally, use the automatic changelog generation button if available.

3. Attach the SBOM file:
    - For stable releases, attach the Software Bill of Materials (SBOM) file to the release by respecting conventions describe in the [SBOM generation](https://handbook.nethserver.org/security/#sbom-software-bill-of-materials) page.
    You can find the SBOM files as artifacts of the GitHub Actions workflow.

## NethVoice release procedure addendum

The following procedure is used by the NethVoice project. It relies on the
`gh ns8-release-module` extension from
[NethServer/gh-ns8-release-module](https://github.com/NethServer/gh-ns8-release-module).

Before using the extension, install it if needed:

```bash
gh extension install NethServer/gh-ns8-release-module
```

1. Run the pre-release checks to verify that all linked issues are in the
   `verified` state:

       gh ns8-release-module check

2. Create the GitHub release and tag, automatically linking closed issues to
   the release notes:

       gh ns8-release-module create --release-name <version> --with-linked-issues

3. Open the newly created GitHub release and review or edit the generated
   release message to ensure it is accurate and complete.

4. If the release includes changes from external components that are not
   tracked as NS8 module repositories, append their current `HEAD` commit
   references to the release notes. This step is temporary and currently
   applies to `nethesis/nethcti-server`. Use the following helper script to
   retrieve the commit SHAs and format them for inclusion:

   ```bash
   #!/bin/bash

   declare -A repositories=(
       ["nethesis/nethcti-server"]="ns8"
   )

   get_last_commit_sha() {
       repo=$1
       branch=$2
       last_commit_sha=$(gh api repos/"$repo"/commits/"$branch" --jq '.sha')
       echo -n "* https://github.com/$repo/commit/$last_commit_sha"
   }

   if ! command -v gh &> /dev/null; then
       echo "GitHub CLI (gh) is not installed. Please install it first."
       exit 1
   fi

   echo "### External components \`HEAD\`"
   for repo in "${!repositories[@]}"; do
       branch=${repositories[$repo]}
       get_last_commit_sha "$repo" "$branch"
       echo
   done
   ```

   Copy the script output and paste it into the release notes.

5. Remove intermediate testing releases that were created during
   the QA cycle:

       gh ns8-release-module clean --release-name <version>

6. Post a comment on each linked issue to notify that the issue
   has been included in the release:

       gh ns8-release-module comment --release-name <version>
