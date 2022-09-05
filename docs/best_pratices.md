---
layout: default
title: Best practices
nav_order: 10
---

# Best practices

* TOC
{:toc}

## Adding more disk space

Most modules run inside `/home` directory.
So, if you're running out of disk space, consider to expand the `/home` filesystem.

On a physical machine, the best option is to setup LVM during the distribution installation.
On the other hand, on a virtual machine you could easily expand the whole root filesystem.

If you can't do any of the above, consider adding a new disk and migrate existing data using below script.
The script will:
- stop all rootless containers
- copy all rootless containers inside the new disk
- reclaim space from the root filesystem
- mount the new disk under `/home`
- restart all rootless containrs

Before running the script, make sure to attach the disk to the machine, format it and mount to a custom
location like `/mnt/temp_disk`.
Then launch the script by passing the mount location as parameter, like:
```
migrate /mnt/temp_disk
```

The `migrate` script:
```bash
#!/bin/bash
set -e

mount_dir=$1

if [ -z "$mount_dir" ]; then
    echo "Please provide the disk mount dir"
    exit 1
fi

if [ ! -d "$mount_dir" ]; then
    echo "Mount dir not found"
    exit 1
fi

if ! which rsync &>/dev/null; then
    echo "rsync not found, please install it before proceed"
    exit 1
fi

device=$(grep $(echo "$mount_dir" | sed 's/\/$//') /etc/mtab | awk '{print $1}')
if [ -z "$device" ]; then
    echo "Device not found for $mount_dir"
    exit 1
fi

for userhome in /home/*[0-9]; do
    moduleid=$(basename $userhome)
    echo "Disabling rootless module ${moduleid}..."
    systemctl stop user@$(id -u $moduleid)
done

echo "Copying files..."
# Copy files
rsync -avrX --delete-after /home/ $mount_dir/

echo "Reclaim space..."
# Delete existing files
for userhome in /home/*[0-9]; do
   rm -rf $userhome
done

echo "Remount disk to /home..."
umount $device
mount $device /home

for userhome in /home/*[0-9]; do
    moduleid=$(basename $userhome)
    echo "Starting rootless module ${moduleid}..."
    systemctl start user@$(id -u $moduleid)
done

echo "Done"
echo

eval $(blkid $device -o export)
echo "Add this line to fstab:"
echo
echo "UUID=$UUID /home $TYPE rw,errors=remount-ro 0 1"

```
