---
layout: default
title: Testing
nav_order: 110
parent: Modules
---

# Testing

NS8 modules use [Robot Framework](https://robotframework.org/) for integration
testing. Tests live in the `tests/` directory of each module repository and run
against a live NS8 cluster. The standard tooling is provided by the
[ns8-github-actions](https://github.com/NethServer/ns8-github-actions) repository,
which includes reusable CI/CD workflows and the `test-ns8-module` script described
below.

The `test-ns8-module` script runs any NS8 module's `tests/` directory with
Robot Framework inside a Podman container, directly from your workstation
against a live NS8 cluster. The venv is cached in a named volume so repeated
runs are fast.

* TOC
{:toc}

## Requirements

- [Podman](https://podman.io/) available in `PATH`
- An SSH private key with access to the target NS8 leader node

## Installation

Download the script, make it executable, and place it in your `PATH`:

<!-- fix curl url //// -->

```bash
curl -o test-ns8-module \
  https://raw.githubusercontent.com/NethServer/ns8-github-actions/refs/heads/v1/scripts/test-module.sh
chmod +x test-ns8-module
sudo mv test-ns8-module /usr/local/bin/
```

## Usage

Enter the NS8 module directory and run the script:

```bash
cd /path/to/ns8-<module>
test-ns8-module <LEADER_NODE> <IMAGE_URL> [robot options...]
```

| Argument | Description |
|---|---|
| `LEADER_NODE` | Hostname or IP of the NS8 leader node |
| `IMAGE_URL` | Container image URL for the module under test |
| `[robot options...]` | Any extra arguments forwarded to the `robot` command |

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `SSH_KEYFILE` | `~/.ssh/id_ecdsa` | Path to the SSH private key |
| `RUN_UI_TESTS` | _(unset)_ | Set to `true` to enable UI/browser tests |

## Examples

Basic run:

```bash
cd ~/git/ns8-mail
test-ns8-module rl1.leader.cluster0.test.org ghcr.io/nethserver/mail:bug-6977
```

Using a custom SSH key:

```bash
cd ~/git/ns8-mail
SSH_KEYFILE=~/.ssh/id_ecdsa test-ns8-module rl1.leader.cluster0.test.org ghcr.io/nethserver/mail:bug-6977
```

With UI tests enabled:

```bash
cd ~/git/ns8-nextcloud
RUN_UI_TESTS=true test-ns8-module rl1.leader.cluster0.test.org ghcr.io/nethserver/nextcloud:latest
```

With UI tests and a custom SSH key:

```bash
cd ~/git/ns8-nextcloud
SSH_KEYFILE=~/.ssh/id_ecdsa RUN_UI_TESTS=true test-ns8-module rl1.leader.cluster0.test.org ghcr.io/nethserver/nextcloud:latest
```

Passing extra Robot Framework options (e.g. run a single test suite):

```bash
cd ~/git/ns8-mail
test-ns8-module rl1.leader.cluster0.test.org ghcr.io/nethserver/mail:latest --suite "Sending mail"
```

Run only tests with a specific tag:

```bash
cd ~/git/ns8-mail
test-ns8-module rl1.leader.cluster0.test.org ghcr.io/nethserver/mail:latest --include smoke
```

Run a single test by name:

```bash
cd ~/git/ns8-mail
test-ns8-module rl1.leader.cluster0.test.org ghcr.io/nethserver/mail:latest --test "Send an email"
```

## How it works

- When `RUN_UI_TESTS=true`, the script uses the [Microsoft Playwright](https://playwright.dev/)
  container image and the `robotframework-browser` library.
- Otherwise, a lightweight `python:3.11-alpine` image is used.
- The Python venv is stored in a named volume (`rftest-cache` or `rftest-cache-ui`).
  It is invalidated automatically when the requirements file checksum changes.
- Robot Framework variables passed to all tests:
  - `NODE_ADDR` — the leader node address
  - `IMAGE_URL` — the module image URL
  - `SSH_KEYFILE` — path to the SSH key inside the container (`/tmp/idssh`)
  - `RUN_UI_TESTS` — whether UI tests are active
- Tests tagged `unstable` are skipped on failure (`--skiponfailure unstable`).
- UI tests must be tagged `ui`; they are excluded automatically when
  `RUN_UI_TESTS` is not `true`.
