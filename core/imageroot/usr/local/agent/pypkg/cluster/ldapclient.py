#
# Copyright (C) 2021 Nethesis S.r.l.
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

import ldap3

def connect(host, port, bind_dn, bind_password, base_dn, **kwargs):
    ldapsrv = ldap3.Server(host, port=int(port))
    ldapconn = ldap3.Connection(ldapsrv,
        user=bind_dn or None,
        password=bind_password or None,
        client_strategy=ldap3.SAFE_SYNC,
        auto_bind=True,
        raise_exceptions=True,
        auto_range=True,
        auto_referrals=False,
    )
    return ldapconn

def get_dn_attributes(conn, dn, lfilter='(objectClass=*)', attributes=[ldap3.ALL_ATTRIBUTES, ldap3.ALL_OPERATIONAL_ATTRIBUTES]):
    response = conn.search(dn, lfilter, ldap3.BASE, attributes=attributes)[2]

    for entry in response:
        if entry['type'] == 'searchResEntry':
            return entry['attributes']

    raise Exception('DN not found')
