---
layout: default
title: Software repositories
nav_order: 7
parent: Core
---

# Software repositories

Module packages are pulled from software repositories. Multiple
repositories can be configured in a NS8 cluster. Each repository contains
the list of available modules
and can be enabled or disabled.
If the repository has the `testing` flag enabled, it will list also non-stable module versions.

Repository metadata are downloaded and cached inside Redis for 3600
seconds.
To expire an existing repository cache, execute: `redis-cli del cluster/repository_cache/<repository_name>`

A NS8 repository must contain a file named `repodata.json` which describes
the content of the repository. The
[ns8-repomd](https://github.com/NethServer/ns8-repomd/) is the default
repository and contains also
the `createrepo.sh` command which creates the metadata.

The `createrepo.sh` command analyzes a directory passed as parameter, if no path is given, it walks the current working directory.
For each directory found, it uses the directory name to query the remote image registry with
[skopeo](https://github.com/containers/skopeo) and retrieve image tags.
If tags are valid [semantic versions](https://semver.org/) they are added to the list of available module version.

Each module that needs to be published inside the repository should have a directory named after the module itself.
The directory should contain:

- a file named `metadata.json`: it contains all module metadata like the name and the URL inside the image registry
- a file named `logo.png`: a PNG file, 256x256 pixels
- a directory named `screenshots`: it can contain one ore more image, each image must be in PNG format,
  with a resolution of 1024x768 pixels

See [dokuwiki directory](https://github.com/NethServer/ns8-repomd/tree/main/dokuwiki) for an example of a complete module
metadata directory.

Execute `createrepo.sh` to generate the `repodata.json`. The format is described [here](https://github.com/NethServer/ns8-scratchpad/blob/main/core/imageroot/var/lib/nethserver/cluster/repodata-schema.json).

