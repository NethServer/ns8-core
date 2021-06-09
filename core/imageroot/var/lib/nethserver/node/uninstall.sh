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
    echo "Deleting rootless module ${moduleid}..."
    systemctl stop user@$(id -u $moduleid)
    loginctl disable-linger "${moduleid}"
    userdel -r "${moduleid}"
done

echo "Stopping the core services"
systemctl disable --now api-server.service redis.service wg-quick@wg0
rm -vf /etc/systemd/system/redis.service.d/wireguard.conf

echo "Wipe Redis DB"
podman volume rm -f redis-data

echo "Uninstalling the core image files"
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
    echo "Deleting rootfull module ${moduleid}..."
    systemctl disable --now agent@${moduleid} && rm -rf "${modulehome}"
done

echo "Some files may be left in the following directories:"
cat ${tmp_dirlist}

systemctl daemon-reload

echo "Wipe Podman storage"
podman system reset -f

echo "Clean up /etc/hosts"
sed -i -e '/ cluster-leader$/ d' -e '/ cluster-localnode$/ d' /etc/hosts
