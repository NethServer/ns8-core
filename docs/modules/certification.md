---
layout: default
title: Community Certification
nav_order: 134
parent: Modules
---

# Application Community Certification

* TOC
{:toc}

## From ideas to applications

Before creating an application for NethServer 8, consider if it's really
necessary. Sometimes, a simple documentation page is sufficient and
requires much less effort. When trying out new ideas, make sure to
document every step and comment in a [community discussion under the
Feature category](http://community.nethserver.org/c/feature).

If the idea gains traction, then it's time to think about a new
application. When creating one, ensure the following requirements are met:

- The application and its components must respect user data privacy and
  integrity. Although this may seem obvious, it's crucial for everyone to
  be aware of any special data treatments performed by the application,
  especially if it relies on third-party cloud services. This allows users
  to evaluate whether those treatments comply with their country's laws
  and personal privacy standards.

- If your application relies on third-party work, make sure it complies
  with their license.

- Create a personal repository on GitHub, for example by starting with the
  [ns8-kickstart
  template](https://github.com/new?template_name=ns8-kickstart&template_owner=NethServer).
  The [How to](../new_module) page provides a brief introduction about
  applications (or modules) development.

- Add a license and copyright notice in the `COPYING` or `LICENSE` file.
  GPL v3+ is highly recommended to ensure your work is distributed and
  used under that license.

- Announce the new application on the [community.nethserver.org
  App](http://community.nethserver.org/c/app) category.

- If desired, submit a pull request for the [NethForge repository
  metadata](https://github.com/NethServer/ns8-nethforge/) to have the
  application listed on the Software Center page. Be sure to include a
  link to the community announcement. If you prefer to setup a personal
  application repository, look at the [Software
  repositories](../../core/software_repositories/) page.

Before being included in NethForge, applications must pass a review
process. The review process leading to inclusion in NethForge is described
in the next section.

## Certification process

The Community Certification ensures that an application is functional,
secure, maintained, and provides sufficient information about the
developer and the software it relies on.

Submitting a pull request (PR) for the [NethForge repository
metadata](https://github.com/NethServer/ns8-nethforge/) starts the
Community certification process for the application.

After approval, at any time, the application can be removed from NethForge
by the NethForge repository maintainers. For example, this can happen if
the current version is no longer maintained and/or has significant security
issues, or if its behavior is deemed malicious.

Here follow some tips that may help to quickly complete the certification
process.

### 1. Metadata and license

0. The PR for NethForge must contain a directory and files as described in
   [Create a
   repository](../../core/software_repositories#create-a-repository).

0. Similarly, fill the [metadata.json](../metadata/) file of the app UI
   with the required information. Provide details for reporting bugs, such
   as a link to the app's bug tracker (you can link to
   `https://github.com/NethServer/dev/issues`), the developer's email, and
   any other relevant contact information.

0. Additionally, in the `metadata.json`, consider linking a document via
   `terms_url` that clarifies any data treatment or application behavior
   the user should be aware of. This is particularly important if the
   application handles sensitive data or must conform to data privacy
   laws. During the certification process, we will pay close attention to
   compliance with EU laws and regulations.

0. Administrator documentation. This is the documentation for the
   administrator who must configure the applications for the final users.
   A public web page link must be provided in the PR description. You can
   decide to write the documentation on your preferred web site or to
   directly create a new [NS8 Admin's
   Manual](https://github.com/NethServer/ns8-docs/) page with a PR. In
   both cases the metadata.json documentation URL will point to the latter
   with an application summary available from there.

0. Integrate the `README.md` file in the code repository with
   documentation for developers and the support team. A detailed
   description of container volumes, environment variables, and special
   maintenance or customization procedures is very important. If something
   is still missing, document it for instance under a "Roadmap",
   "Limitations", "TODO" section. Provide also a link to your Community
   profile to contact you.

0. The application and its parts must have no licensing conflicts.
   Licenses must permit unlimited distribution of the application.

### 2. Platform compliancy

0. The implementation of platform-defined behaviors, where present, must work
   correctly. Pay attention to:

    - Application Creation and Removal
    - Clone and Move
    - Backup and Restore
    - Change of cluster Email notification settings
    - HTTP routes management
    - TLS certificates management
    - Bind to one or more User Domains
    - Node firewall management
    - Node ports allocation

0. Applications must avoid requiring root privileges (rootfull) and should
   operate with the minimum necessary authorizations. Exceptions may be
   requested but must be justified and approved on a case-by-case basis.

0. The UI code respects Internationalization (i18n) and [NS8 Weblate
   project](https://hosted.weblate.org/projects/ns8/) has a component for
   your app Localization (l10n). The app translation license must be
   compatible with the Weblate Libre Hosted plan. For this purpose, if
   your app license is not compliant, choose a translation license from
   [choosealicense.com](https://choosealicense.com/).

### 3. Application updates

Application updates must obey the same rules of the core, described in the
[Development process](../../development_process#update-rules).
