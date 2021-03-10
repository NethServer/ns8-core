#!/bin/bash

#
# Clean up data and control plane init.sh
#

shopt -s nullglob

echo "[NOTICE] Stop control plane agents"
systemctl stop node-agent.service module-agent@\*.service

echo "[NOTICE] Uninstalling the control plane"
find /etc/systemd/system/  \( -name module\*.service -o -name node-agent.service \) -delete

rm -rvf /usr/local/share/{cplane,dplane,agent}\
    /etc/sysctl.d/80-nethserver.conf \
    /etc/systemd/user/module-agent.service /etc/systemd/user/module-init.service \
    /usr/local/etc/node-agent.env

for userhome in /home/*; do
    moduleid=$(basename $userhome)
    echo "[NOTICE] Deleting rootless module ${moduleid}..."
    loginctl disable-linger "${moduleid}"
    userdel -r "${moduleid}"
done

for modulehome in /var/local/*; do
    moduleid=$(basename $modulehome)
    echo "[NOTICE] Stopping rootfull module ${moduleid}..."
    systemctl disable --now "${moduleid}"
    rm -f "/etc/systemd/system/${moduleid}.service"
    echo "[NOTICE] Deleting rootfull module ${moduleid}..."
    rm -rvf "${modulehome}" /usr/local/etc/${moduleid}.env
done

systemctl daemon-reload