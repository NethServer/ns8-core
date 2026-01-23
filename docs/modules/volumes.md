---
layout: default
title: Volumes
nav_order: 30
parent: Modules
---

# Volumes

To enable data persistence and to obtain better performance, application
data must reside in Podman volumes. A volume, passed via `podman run
--volume` parameter can be of these types:

- Local directory, bind-mounted into the container filesystem.
- Single local file, copied into the container filesystem when podman run
  is invoked.
- Podman named volume.

The following example shows how multiple volumes can be defined for a
service in the Systemd `dovecot.service` unit:

```text
[Service]
...
ExecStart=/usr/bin/podman run \
    ...
    --volume=dovecot-data:/var/lib/vmail \
    --volume=./tls-certs:/etc/ssl/dovecot:z \
    --volume=dovecot-dict:/var/lib/dovecot/dict:z \
    --volume=dovecot-lmtp:/var/lib/umail:z \
    --volume=dovecot-custom:/etc/dovecot/local.conf.d:z \
    ${MAIL_DOVECOT_IMAGE}
```

The above `podman run` arguments define:

* Four named volumes: `dovecot-data`, `dovecot-dict`, `dovecot-lmtp`, and
  `dovecot-custom`. Named volumes are typically used to store application
  data persistently. In some cases, they are shared between containers—for
  example, to enable communication through a Unix-domain socket.

* One bind-mounted directory: `./tls-certs`. This directory, located under
  the current working directory when `podman run` is invoked, is
  bind-mounted to the container’s `/etc/ssl/dovecot`. This approach is
  often used to pass configuration files—or TLS certificates, in this
  case—to the container. A third option, single-file volumes, is described
  in the Podman documentation.

Note the `:z` flag, present in some volume definitions: it configures the
SELinux volume label in "shared" mode (type `container_file_t`), which is
the recommended choice for rootless containers. Since the label is
recursively checked when the container is started, the volume
`dovecot-data` that may contain many files is exempted to avoid long
startup times. If you copy files from external sources into a volume,
ensure the label type is `container_file_t` to avoid file access errors.
Refer to `podman run` man page for differences between volume flags.

Inspect Podman named volumes with:

    podman volume ls

## Suggest volume assignment to additional disks

When an application is created in NS8, or when a service is started for
the first time, named volumes are created under a Podman predefined
directory. You can check its path with:

    {% raw %}podman system info --format='{{.Store.VolumePath}}'{% endraw %}

The path depends on whether the command is executed by root (for rootful
modules) or by a normal user (for rootless modules).

The default path resides under the root filesystem. It is possible to mark
one or more application volumes to be presented for optional assignment to
additional disks, present in the NS8 installation node.

The `org.nethserver.volumes` [image label](../images/#image-labels) is managed by the
application developer. Define this application label like:

     org.nethserver.volumes = dovecot-data

When the Mail application is installed in a NS8 cluster, this label
triggers an UI modal window asking the Sysadmin where to store its data,
provided the destination node has an additional disk already configured.
If not, the volume is created in the default disk.

The label value syntax is a list of words separated by spaces. Each word
must correspond to a named volume. List only the volumes you consider
suitable for additional disks. 

Assume that additional disks are **slower than the root disk**, but have
**bigger size**.

## Configure volume assignments for a node

The system administrator can configure a NS8 node to assign specific
application volumes to custom base directories. These assignments are
stored in `/etc/nethserver/volumes.conf` and evaluated during module
creation.

Assignments defined in `volumes.conf` can extend beyond the volume names
declared by the app developer with the `org.nethserver.volumes` image
label explained in the previous section. However, settings collected from
the UI take precedence over those in `volumes.conf`.

The `volumes.conf` file uses a Windows-INI–like format, where section names
correspond to applications and key-value pairs define named volume
assignments with their system mount points. It can be edited manually with
a text editor or managed through the `volumectl` command.

This is a brief list of `volumectl` invocations for different scenarios.

1. List available base directories (filesystem mountpoints). For each
   directory the mounted disk label, size and used space are printed.

       # volumectl list-base-paths
       /srv/disk00 (MYDISK) size=2.0G used=46.6M

2. Assign base path for `dovecot-data` volume of `mail` applications.

       # volumectl add-volume dovecot-data --target /mnt/disk00 --for mail
       Added volume target for [mail]: dovecot-data → /mnt/disk00

3. Similar to the previous example, but limits the assignment only to the
   next Nextcloud installation.

       # volumectl add-volume nextcloud-app-data --target /mnt/disk01 --for nextcloud --next-only
       Added volume target for [nextcloud4]: nextcloud-app-data → /mnt/disk01
       # volumectl add-volume nextcloud-app-data --target /mnt/disk00 --for nextcloud
       Added volume target for [nextcloud]: nextcloud-app-data → /mnt/disk00

4. Remove an existing volume assignment.

       # volumectl remove-volume nextcloud-app-data --for nextcloud
       (command has no output on success)

Look at `/etc/nethserver/volumes.conf` to check the current settings:

       cat /etc/nethserver/volumes.conf

Output example:

```text
[mail]
dovecot-data = /mnt/disk00

[nextcloud4]
nextcloud-app-data = /mnt/disk01

[nextcloud]
nextcloud-app-data = /mnt/disk00
```

This configuration file example shows how specific volume targets apply to
different applications. It will store `nextcloud-app-data` data
volume of `nextcloud4` application under `/mnt/disk01`, whilst data of
other Nextcloud applications will reside under `/mnt/disk00`.
