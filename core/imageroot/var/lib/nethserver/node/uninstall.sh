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
    systemctl stop user@$(id -u $moduleid)
    loginctl disable-linger "${moduleid}"
    userdel -r "${moduleid}"
done

echo "[NOTICE] Stopping the core services"
systemctl disable --now api-server.service redis.service wg-quick@wg0

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

for modulehome in /var/lib/nethserver/*; do
    if [[ ! -d ${modulehome} ]]; then
      continue
    fi
    moduleid=$(basename $modulehome)
    echo "[NOTICE] Deleting rootfull module ${moduleid}..."
    systemctl disable --now agent@${moduleid} && rm -rf "${modulehome}"
done

echo "[NOTICE] Some files may be left in the following directories:"
cat ${tmp_dirlist}

systemctl daemon-reload

echo "[NOTICE] Wipe Podman storage"
podman system reset -f "${@}"
