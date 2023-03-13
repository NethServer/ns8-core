#!/bin/bash

#
# Copyright (C) 2022 Nethesis S.r.l.
# http://www.nethesis.it - nethserver@nethesis.it
#
# This script is part of NethServer.
#
# NethServer is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License,
# or any later version.
#
# NethServer is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with NethServer.  If not, see COPYING.
#

#
# Erase data and uninstall the core scripts
# Some files can be still left in place: look at this
# script output for their paths.
#

shopt -s nullglob

tmp_dirlist=$(mktemp)
tmp_srclist=$(mktemp)

cat /var/lib/nethserver/node/state/coreimage.lst > ${tmp_srclist}

trap "rm -f ${tmp_dirlist} ${tmp_srclist}" EXIT

for userhome in /home/*[0-9]; do
    moduleid=$(basename $userhome)
    echo "Deleting rootless module ${moduleid}..."
    loginctl disable-linger "${moduleid}"
    systemctl stop user@$(id -u $moduleid)
    loginctl terminate-user "${moduleid}"
    userdel -r "${moduleid}"
done

echo "Clean up firewalld core rules"
firewall-cmd --permanent --remove-service=http --remove-service=https --remove-service=ns-wireguard >/dev/null
firewall-cmd --permanent --delete-service=ns-wireguard

wg0_cluster_network=$(runagent redis-get cluster/network)
if [[ -n "${wg0_cluster_network}" ]]; then
  firewall-cmd --permanent --zone=trusted --remove-source="${wg0_cluster_network}" >/dev/null
fi

firewall-cmd --reload

echo "Stopping the core services"
systemctl disable --now api-server.service redis.service wg-quick@wg0.service phonehome.timer
rm -vf /etc/systemd/system/redis.service.d/wireguard.conf \
   '/etc/systemd/system/wg-quick@wg0.service.d/persist-wg0-config.conf' \
   /etc/wireguard/wg0.conf
userdel -r api-server
ip link delete wg0 # force removal of wg0 interface

echo "Wipe Redis DB"
podman volume rm -f redis-data

for modulehome in /var/lib/nethserver/*[0-9]; do
    if [[ ! -d ${modulehome} ]]; then
      continue
    fi
    moduleid=$(basename $modulehome)
    echo "Deleting rootfull module ${moduleid}..."
    readarray -t units < <(find /etc/systemd/system -type f -a \( \
      -name "${moduleid}*.service" \
      -o -name "backup*-${moduleid}.service" \
      -o -name "backup*-${moduleid}.timer" \) \
    )
    # shellcheck disable=SC2046
    systemctl disable --now "agent@${moduleid}" $(basename -a "${units[@]}")
    rm -rfv "${modulehome}" "${units[@]}"
done

echo "Deleting cluster and agent core modules"
systemctl disable --now \
  agent@node.service \
  agent@cluster.service
rm -rf /var/lib/nethserver/cluster /var/lib/nethserver/node

echo "Removing main directories"
rm -rf /usr/local/agent /var/lib/nethserver /etc/nethserver

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
) <${tmp_srclist}

echo "Some files may be left in the following directories:"
cat ${tmp_dirlist}

systemctl daemon-reload

echo "Wipe Podman storage"
podman system reset -f

echo "Clean up /etc/hosts"
sed -i -e '/ cluster-leader$/ d' -e '/ cluster-localnode$/ d' /etc/hosts
