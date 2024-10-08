---
layout: default
title: Software repositories
nav_order: 7
parent: Core
---

# Software repositories

In NS8, applications are available from remote application indexes known as "repositories."

## Cluster configuration

During installation, two repositories are configured:

- `default` (enabled)
- `nethforge` (disabled)

Key points:

- You can configure additional repositories in an NS8 cluster.
- Each repository can be in either an **enabled** or **disabled** state.
- The list of available modules in NS8 is generated by merging data from
  all enabled repositories.
- If two repositories contain modules with the same name, the module from
  the repository with higher priority is used, and others are ignored.
- Repository priority is determined by the alphabetical order of
  repository names, with names later in the alphabet (e.g., "Z") having
  higher priority than those earlier (e.g., "A").
- If a repository has the **testing** flag enabled, pre-releases are
  treated as stable and are included in application update calculations.

Repository configurations are stored in a Redis HASH with the key prefix
`cluster/repository/`. For example, to view the configuration of the
`default` repository, run:

```bash
redis-cli HGETALL cluster/repository/default
```

To enable the testing flag (for QA purposes only), use the following command:

```bash
redis-cli HSET cluster/repository/default testing 1
```

Note: If the testing flag is manually enabled, it will be reset to `0`
after the next core update.

Repository metadata is downloaded and cached in Redis for 3600 seconds. To
clear an existing repository cache, run:

```bash
redis-cli del cluster/repository_cache/<repository_name>
```

## Create a repository

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

Execute `createrepo.sh` to generate the `repodata.json`. The format is described [here](https://github.com/NethServer/ns8-core/blob/main/core/imageroot/var/lib/nethserver/cluster/repodata-schema.json).

