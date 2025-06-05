#!/bin/bash

#
# Copyright (C) 2025 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#
#
# This script migrates the home directory to a new mounted disk
#
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
rsync -avrX --delete-after /home/ $mount_dir/
echo "Reclaim space..."
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
echo "UUID=$UUID /home $TYPE defaults 0 1"
