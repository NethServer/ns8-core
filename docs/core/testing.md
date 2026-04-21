---
layout: default
title: Testing suite
nav_order: 11
parent: Core
---

# Testing suite

* TOC
{:toc}

## Requisites

The tests have been created using [Robot Framework](https://robotframework.org/) and the execution happens inside a Podman container. The only dependency required
in the host system is [Podman](https://podman.io/getting-started/installation).

## Tests structure

All tests are in the `tests` directory of NS8 core, in the root of the `tests` directory the file `pythonreq.txt`
specifies the dependencies required and each sub-directory contains a test suite with the related tests.

The test suite and the related tests are executed in alphabetical order, the `__init__.robot` can be used for the suite
initialization. More details in the [Robot Framework documentation](https://robotframework.org/robotframework/latest/RobotFrameworkUserGuide.html#files-and-directories).

```
core/tests/
├── 00_<test_suite>
│   ├── 00_<test>.robot
│   ├── 10_<test>.robot
│   ├── 20_<test>.robot
│   ├── ...
│   ├── ...
│   └── __init__.robot
├── 01_<test_suite>
│   ├── 00_<test>.robot
│   ├── 10_<test>.robot
│   ├── 20_<test>.robot
│   ├── ...
│   ├── ...
│   └── __init__.robot
├── ...
│   ├── ...
│   ├── ...
│   └── __init__.robot
├── pythonreq.txt
└── __init__.robot
```

## Tests execution

NS8 core uses Robot Framework for integration testing. Tests run inside a Podman container against a machine of your choice using the `test-ns8-module` script. The standard tooling is provided by the [ns8-github-actions](https://github.com/NethServer/ns8-github-actions) repository.

### Requirements

- [Podman](https://podman.io/) available in `PATH`
- An SSH private key with access to the target NS8 leader node

### Installation

Download the script, make it executable, and place it in your `PATH`:

<!-- fix curl url //// -->

```bash
curl -o test-ns8-module \
  https://raw.githubusercontent.com/NethServer/ns8-github-actions/refs/heads/v1/scripts/test-module.sh
chmod +x test-ns8-module
sudo mv test-ns8-module /usr/local/bin/
```

### Usage

Enter the core directory, then run tests:

```bash
cd core/
test-ns8-module <LEADER_NODE> [robot_options...]
```

| Argument | Description |
|---|---|
| `<LEADER_NODE>` | Hostname or IP of the NS8 leader node |
| `[robot_options...]` | Extra arguments forwarded to the `robot` command (e.g. `--include`, `--exclude`) |

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `SSH_KEYFILE` | `~/.ssh/id_ecdsa` | SSH private key to use for cluster connection |
| `RUN_UI_TESTS` | _(unset)_ | Set to `true` to enable UI/browser test cases |
| `COREMODULES` | _(unset)_ | Space-separated or comma-separated list of core modules to install during cluster setup |

### Test tags

Test cases can be marked with [Robot Framework tags](https://robotframework.org/robotframework/latest/RobotFrameworkUserGuide.html#tagging-test-cases) to include or exclude tests during execution:

- `backend` — All tests in `10__cluster_sanity` suite
- `ui` — UI test cases (excluded by default, enable with `RUN_UI_TESTS=true`)
- `install` — Tests related to NS8 installation and cluster creation
- `uninstall` — Tests related to NS8 removal
- `unstable` — Tests automatically skipped on failure

To skip installation and uninstallation tests (useful for testing an already-installed cluster):

```bash
test-ns8-module <LEADER_NODE> --exclude install --exclude uninstall
```

To enable UI tests:

```bash
RUN_UI_TESTS=true test-ns8-module <LEADER_NODE>
```

### Examples

Basic run (substitute with your leader node hostname):

```bash
test-ns8-module rl1.leader.cluster0.test.nethserver.org
```

Using a custom SSH key:

```bash
SSH_KEYFILE=~/.ssh/id_ecdsa test-ns8-module rl1.leader.cluster0.test.nethserver.org
```

With UI tests enabled:

```bash
RUN_UI_TESTS=true test-ns8-module rl1.leader.cluster0.test.nethserver.org
```

With specific core modules:

```bash
COREMODULES="ghcr.io/nethserver/core:latest ghcr.io/nethserver/traefik:feat-7544" test-ns8-module rl1.leader.cluster0.test.nethserver.org
```

Skipping installation and uninstallation tests:

```bash
test-ns8-module rl1.leader.cluster0.test.nethserver.org --exclude install --exclude uninstall
```

## Testing environment

A Terraform configuration is available in the [`ns8-terraform-infra`](https://github.com/NethServer/ns8-terraform-infra) repository to provision a clean infrastructure for test execution.

### Setup

1. Clone the repository and follow its setup instructions
2. [Export the generated SSH keys](https://github.com/NethServer/ns8-terraform-infra#default-ssh-keys-pair)
3. Set `SSH_KEYFILE` to point to the generated key

### Usage

Running tests with the Terraform-generated key:

```bash
SSH_KEYFILE=../../ns8-terraform-infra/key test-ns8-module <LEADER_NODE>
```

Accessing nodes with the generated key:

```bash
ssh -i ../../ns8-terraform-infra/key <leader_node>
```

### Example

For a domain named `test.nethserver.org` running on Rocky Linux 9 in the `cluster0` workspace:

```bash
SSH_KEYFILE=../../ns8-terraform-infra/key test-ns8-module rl1.leader.cluster0.test.nethserver.org
```
