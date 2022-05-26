---
layout: default
title: FAQ
nav_order: 9
---

## FAQ: Frequently Asked Questions

* TOC
{:toc}

### What is NethServer 8 (NS8)?

NethServer 8, aka NS8, is a simple-to-use container orchestrator.
It's aimed for Linux sysadmins who need the flexibility of containers,
but do not want to mess around with a complex tools like Kubernetes.

Like NethServer 7, NS8 hides the system complexity behind a ready-to-use
beautiful Web User Interface.

Do you want to easily manage applications without worrying about installation details?
NS8 is for you.

### What is the status of NS8?

NethServer 8 is a work in progress, currently in a pre-release status.  Answers in this FAQ describing its capabilities describe the goals of the project, not necessarily its current state.

### What is the difference between NethServer 7 (NS7) and NethServer 8?

NethServer 7 is an operating system built on top of CentOS 7.
It is heavily-coupled with the underlying OS. It installs applications
using RPMs package and configure them with a configuration management
system called e-smith.

NethServer 8 can run on different [Linux distributions]({{site.baseurl}}/quickstart#system-requirements).
It installs applications using containers and configure them using
a brand new architecture which supports also a multi-node configuration.

As NS7, NS8 is perfectly suited for small offices and medium enterprises.

### What can I do with NethServer 8?

You can serve an Active Directory domain with Samba 4, or just a simple OpenLDAP server.
You can host multiple instances of the same application connected to different local or
remote domains.

If you need a new domain, just create a new node and join it to the existing cluster:
you can spread the workload across multiple machines.
Of course, everything is under your control just from the UI of the main node.
No matter if nodes are local or remotes, they are securely connected with an encrypted VPN.

If you like an all-in-one server. you can do it too.

### Can I migrate from NethServer 7 to NethServer 8?

Migrating to the next major release is a primary goal, as always has been.
The migration procedure is still not ready but developers will work hard to make it *safe* and *painless*.

### How can I install NethServer 8?

Take a look to the [Quickstart]({{site.baseurl}}/quickstart#core-installation).

### Where can I download it?

No need to download it. Just install on a [supported distribution]({{site.baseurl}}/quickstart#system-requirements).

### Can I install NethServer on a virtual machine?

Yes. You can choose any hypervisor or cloud provider supporting [these distributions]({{site.baseurl}}/quickstart#system-requirements).

### Can I install it on bare metal?

Sure. See the question below.

### What hardware is supported?

Anything supported by [these distributions]({{site.baseurl}}/quickstart#system-requirements),

### Does NethServer 8 have Cockpit?

No, NethServer 8 has a brand new [User Interface]({{site.baseurl}}/ui).

### Can I use NethServer 8 as UTM firewall?

No, NethServer 8 doesn't have any firewall capabilities.

