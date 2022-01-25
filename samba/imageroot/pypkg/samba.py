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

import ipaddress as ipm
import subprocess
import socket

class SambaException(Exception):
    pass

class IpNotPrivate(SambaException):
    pass

class IpNotAvailable(SambaException):
    pass

class IpBindError(SambaException):
    def __init__(self, ipaddr, message):
        self.ipaddr = ipaddr
        self.message = message
        super().__init__(self.message)

def ipaddress_check(ipaddress):

    addr = ipm.ip_address(ipaddress)
    # See Python docs: https://docs.python.org/3.9/library/ipaddress.html#ip-addresses
    if not addr.is_private or addr.is_unspecified or addr.is_reserved or addr.is_loopback or addr.is_link_local:
        raise IpNotPrivate()

    ipproc = subprocess.run(
        ["ip", "-o", "address", "show", "to", f"{ipaddress}/32", "up"],
        capture_output=True,
        text=True,
    )

    if not ipaddress in ipproc.stdout:
        raise IpNotAvailable()

    for xipaddr in [ipaddress, '127.0.0.1']:
        for tcp_port in [53, 88, 636, 464, 445, 3268, 3269, 389, 135, 139]:
            try:
                with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sk:
                    sk.bind((xipaddr, tcp_port))
            except Exception as ex:
                raise IpBindError(xipaddr, f"Address {xipaddr}:{tcp_port} bind failed: {ex}") from ex

    return True