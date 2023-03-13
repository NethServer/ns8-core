#
# Copyright (C) 2023 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import os
import socket

def initialize_wgconf(ipaddr, listen_port=55820, peer=None):
    """Initialize wg0.conf to bootstrap the VPN link. The file is then
       overwritten by the `wg-quick save wg0` by the Systemd unit stop
       handler and by the `apply-vpn-routes` command"""

    with open('/etc/nethserver/wg0.key') as fpkey:
        private_key = fpkey.read().strip()

    with open('/etc/nethserver/wg0.pub') as fpkey:
        public_key = fpkey.read().strip()

    oldmask = os.umask(0o77)
    # Create a new file beside our target file path:
    with open('/etc/wireguard/wg0.conf.new', 'w') as wgconf:
        os.umask(oldmask)

        # Write the Interface head section:
        wgconf.write(f"[Interface]\n")
        wgconf.write(f"Address = {ipaddr}\n")
        wgconf.write(f"ListenPort = {listen_port}\n")
        wgconf.write(f"PrivateKey = {private_key}\n\n")

        # Append Peer section for VPN bootstrap:
        if peer:
            peer_hostname, peer_port = peer['endpoint'].rsplit(':', 1)

            # Resolve the address to an IP string:
            peer_ep_address = socket.getaddrinfo(peer_hostname, peer_port, proto=socket.IPPROTO_UDP)[0][4][0]
            if ':' in peer_ep_address:
                # for IPv6, wrap in square brakets "[]"
                peer_ep_address = f'[{peer_ep_address}]'

            wgconf.write(f"[Peer]\n")
            wgconf.write(f"PublicKey = {peer['public_key']}\n")
            wgconf.write(f"AllowedIPs = {peer['ip_address']}\n")
            wgconf.write(f"PersistentKeepalive = 25\n")
            wgconf.write(f"Endpoint = {peer_ep_address}:{peer_port}\n")

    # Overwrite the target file path:
    os.rename('/etc/wireguard/wg0.conf.new', '/etc/wireguard/wg0.conf')
