#!/bin/bash

#
# Erase data and uninstall the core scripts
# Some files can be still left in place: look at this
# script output for their paths.
#

shopt -s nullglob

tmp_dirlist=$(mktemp)
trap "rm -f ${tmp_dirlist}" EXIT

for userhome in /home/*; do
    moduleid=$(basename $userhome)
    echo "[NOTICE] Deleting rootless module ${moduleid}..."
    loginctl disable-linger "${moduleid}"
    while loginctl show-user "${moduleid}" &>/dev/null; do
        sleep 1
    done

    userdel -r "${moduleid}"
done

echo "[NOTICE] Stopping the core services"
systemctl disable --now api-server.service agent@{cluster,node}.service redis.service

echo "[NOTICE] Uninstalling the core image files"
(
  while read image_entry; do
    # Add a leading / to image_entry
    image_entry="/${image_entry}"
    if [[ -d ${image_entry} ]]; then
        echo ${image_entry} >> $tmp_dirlist
        continue
    fi
    
    [[ -e ${image_entry} ]] && rm -vf ${image_entry}
  done
) </var/lib/nethserver/node/state/image.lst

echo "[NOTICE] Disable WireGuard wg0 interface"
systemctl disable --now wg-quick@wg0

echo "[NOTICE] Some files may be left in the following directories:"
cat ${tmp_dirlist}

systemctl daemon-reload

echo "[NOTICE] Wipe Podman storage"
podman system reset -f "${@}"
