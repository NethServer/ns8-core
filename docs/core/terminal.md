---
layout: default
title: Terminal
nav_order: 18
parent: Core
---

# Terminal

Cluster administrators can open an interactive shell on any cluster node from
cluster-admin. The pseudo-terminal is created by `sshd` on the target node;
api-server opens the SSH connection over the WireGuard VPN and relays the byte
stream to the browser.

* TOC
{:toc}

## What the feature is

**A jump host.** The `open-terminal` grant does not give root. It gives the right
to reach port 22 of a node through the cluster VPN, which is usually not routable
from an administrator's workstation. Authentication is still performed by `sshd`
with a system account and password supplied by the user.

## Two administrator populations

The design assumes that "NS8 cluster administrator" and "operating system
administrator" may be two different people.

`open-terminal` can be granted on a single node to an operator who holds no other
privilege there, and conversely an NS8 `owner` cannot obtain a shell without
system credentials. This is why the authentication factor always comes from the
user and never from NS8, and why nothing in NS8 writes to `authorized_keys`.

A key enrolment action was considered and rejected: `owner` holds the `*` action
pattern on every node, so it would hold any enrolment action by construction,
enrol itself and obtain root without a system factor. "Grantable to `owner` only"
is not a restriction, it is the maximal grant.

## Enabling the terminal on a node

The terminal is disabled on every node by default, and enabling it **modifies the
node's sshd configuration**.

`enable-node-terminal` sets the `node/<id>/terminal` flag, then asks the node to
write `/etc/ssh/sshd_config.d/90-ns8-terminal.conf`:

```
Match Address <cluster/network>
    PermitRootLogin yes
    PasswordAuthentication yes
    KbdInteractiveAuthentication yes
Match all
```

`prohibit-password` has been OpenSSH's compiled-in default since 7.0, and it
refuses both password and keyboard-interactive authentication for root. Without
this block, root cannot log in with a password on either supported distribution
family.

The block grants that access **from the cluster VPN only** — never from the LAN,
never from the Internet — and only on a node whose flag is set. It does not
weaken the two-population rule: it grants the right to *present* a system factor,
not access without one.

On a node hardened with `PasswordAuthentication no`, this block re-allows password
authentication from the VPN. The settings switch states this.

`disable-node-terminal` clears the flag and removes the file.

### Address handling

`cluster/network` holds a single CIDR. Do not build the `Match Address` value
from api-server's `getClusterNetworks()`, which prepends the loopback addresses:
the support tunnel DNATs port 22 to `127.0.0.1`, so loopback here would open root
password login on a path that is public-key only today.

### Convergence

`set-terminal-sshd` takes no parameter. It reads `node/<id>/terminal` and
converges the drop-in to it, so a caller can only realign the node with the
cluster policy, never open password authentication on its own — which matters
because `owner` holds that action anyway.

The flag is written first in both directions. The worst transient state is "flag
set, drop-in not yet written": the handshake fails, visibly, and nothing is open.
The opposite order risks leaving `sshd` accepting root passwords while the
interface reports the terminal as disabled.

Three triggers, fastest first:

1. `enable-node-terminal` and `disable-node-terminal` push the node action right
   after writing the flag;
2. the `node/<id>/tasks` queue is durable, so a powered-off node converges when it
   comes back, with no timer;
3. the `update-core.d/55terminal_sshd` hook realigns. This is the long-stop, not
   the main path: weeks can pass between two core updates.

`sshd -t` runs before any reload, and the drop-in is reverted if validation
fails: a syntax error would otherwise take `sshd` down on that node.

The action fails loudly when `sshd_config` has no
`Include /etc/ssh/sshd_config.d/*.conf` directive, because the drop-in would be
inert and reporting success would be a lie.

## Preconditions

`probe-terminal-access` reports whether a session can succeed, and the browser
runs it before asking for credentials — so a user is not made to type a root
password for a handshake that cannot work.

It returns three booleans and the port, and nothing else: task output is readable
by any authenticated user because GET requests bypass authorization, so the raw
`sshd` configuration must not be exposed.

The probe runs `sshd -T -C user=root,addr=…,laddr=…,lport=…`. Without `-C`,
`sshd` ignores every `Match` block, including the one the opt-in installs, so the
answer would always be the pre-activation policy.

Grant `probe-terminal-access` wherever `open-terminal` is granted:

```
grant-actions --action open-terminal --on node/2 --to <role>
grant-actions --action probe-terminal-access --on node/2 --to <role>
```

Note that `open-terminal` is an authorization name, not an action directory:
nothing executes it. It is the value api-server matches against the grants when
a session is requested, so it will not appear in `api-cli list-actions`. The
action name must stay outside the `get-*`, `list-*`, `show-*` and `read-*`
patterns, which the built-in `reader` role holds.

## Host keys

The same probe publishes the host keys `sshd` actually serves into
`node/<id>/ssh`. api-server accepts those keys and refuses to connect when none
is published.

This is **not** pinning: a legitimately regenerated key is accepted as soon as it
is published, and so would a substituted key if an attacker controlled the node —
a case where the node is lost anyway. The property is that the key arrives over
the cluster's authenticated channel, Redis ACL plus WireGuard, which is stronger
than trust-on-first-use.

`HostKeyAlgorithms` is set from the published key types, otherwise a server
offering RSA before a published Ed25519 key produces a phantom refusal.

## Auditing

`audit.db` records session open and close, failed handshakes, the cluster-admin
identity, the SSH username, the browser source address and the **local TCP port**
of the SSH connection. `sshd` logs `Accepted … port N` and `Failed password …
port N`, so that port makes the join with the node journal exact rather than
merely chronological.

Node-side attribution is otherwise degraded: `last` and `lastb` show the leader's
address for every session. As a complement, api-server runs `logger` over a
second SSH channel before opening the pseudo-terminal, so the node journal carries
the cluster-admin identity — the operating system administrator has no access to
`audit.db`.

Failed handshakes must be audited: without them, a cluster administrator guessing
root passwords is invisible both in `audit.db` and in `lastb`.

## What this does not protect

api-server relays the stream, so it sees everything typed, including the system
password. That is inherent to any browser terminal: even without a login prompt
of ours, users type `sudo`, `mysql -p` or a nested `ssh` inside the shell.
End-to-end encryption would not change it, since api-server serves the JavaScript
and could substitute the key material without the browser having any anchor to
notice.

The `open-terminal` grant is also a password oracle: it allows testing system
passwords against `sshd` from an address that sits in firewalld's `trusted` zone,
bypassing the network restrictions an administrator may have set for external
access. Resistance rests **entirely** on api-server's throttle — `pam_faillock`
has `even_deny_root` disabled by default, so root is not locked out node-side, and
faillock is not in Debian's default stack. The throttle is counted per node and
SSH account, survives WebSocket reconnections, and grows its delay.

The per-node flag is a policy switch, not a security boundary. The authoritative
boundary is `sshd` and the system password.

Cluster backups list keys explicitly, so `node/<id>/terminal` is not saved and a
restore lands with the terminal disabled and no drop-in.
