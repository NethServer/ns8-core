#
# Copyright (C) 2023 Nethesis S.r.l.
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
from .exceptions import LdapclientEntryNotFound
from .base import LdapclientBase

class LdapclientRfc2307(LdapclientBase):

    def get_group(self, group):
        # Escape group string to build the filter assertion:
        escgroup = ldap3.utils.conv.escape_filter_chars(group)
        response = self.ldapconn.search(self.base_dn, f'(&(objectClass=posixGroup)(cn={escgroup}){self._get_groups_search_filter_clause()})',
            attributes=['cn', 'memberUid', 'description'],
        )[2]

        def lget_user(uid):
            # Escape uid string to build the filter assertion:
            escuid = ldap3.utils.conv.escape_filter_chars(uid)
            user = self.ldapconn.search(self.base_dn, f'(&(objectClass=posixAccount)(objectClass=inetOrgPerson)(uid={escuid}){self._get_users_search_filter_clause()})',
                attributes=['displayName', 'uid'],
            )[2][0]['attributes']

            return {
                'user': user['uid'][0],
                'display_name': user.get('displayName') or "",
            }

        for entry in response:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals
            return {
                "group": entry['attributes']['cn'][0],
                "description": __class__.get_entry_description(entry),
                "users": [lget_user(uid) for uid in  entry['attributes']['memberUid']],
            }

        raise LdapclientEntryNotFound()

    def list_groups(self):
        response = self.ldapconn.search(self.base_dn, f'(&(objectClass=posixGroup){self._get_groups_search_filter_clause()})',
            attributes=['cn', 'memberUid', 'description'],
        )[2]

        groups = []
        for entry in response:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals
            groups.append({
                "group": entry['attributes']['cn'][0],
                "description": __class__.get_entry_description(entry),
                "users": entry['attributes']['memberUid'],
            })
        return groups

    def get_user(self, user):
        entry = self.get_user_entry(user)

        def get_memberof(user):
            # Escape user string to build the filter assertion:
            escuser = ldap3.utils.conv.escape_filter_chars(user)
            groups = []
            gresponse = self.ldapconn.search(self.base_dn, f'(&(objectClass=posixGroup)(memberUid={escuser}){self._get_groups_search_filter_clause()})',
                attributes=['cn', 'memberUid', 'description'],
            )[2]
            for gentry in gresponse:
                if gentry['type'] != 'searchResEntry':
                    continue # ignore referrals
                groups.append({
                    "group": gentry['attributes']['cn'][0],
                    "description": __class__.get_entry_description(gentry),
                })
            return groups

        return {
            "user": entry['attributes']['uid'][0],
            "display_name": entry['attributes'].get('displayName') or "",
            "groups": get_memberof(user),
            "locked": entry['attributes'].get('pwdAccountLockedTime', []) != [],
        }

    def get_user_entry(self, user, lextra_attributes=[]):
        # Escape user string to build the filter assertion:
        escuser = ldap3.utils.conv.escape_filter_chars(user)
        response = self.ldapconn.search(self.base_dn, f'(&(objectClass=posixAccount)(objectClass=inetOrgPerson)(uid={escuser}){self._get_users_search_filter_clause()})',
            attributes=['displayName', 'uid'] + self.filter_schema_attributes(['pwdAccountLockedTime'] + lextra_attributes),
        )[2]

        for entry in response:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals

            return entry

        raise LdapclientEntryNotFound()

    def list_users(self):
        response = self.ldapconn.search(self.base_dn, f'(&(objectClass=posixAccount)(objectClass=inetOrgPerson){self._get_users_search_filter_clause()})',
            attributes=['displayName', 'uid'] + self.filter_schema_attributes(['pwdAccountLockedTime']),
        )[2]

        users = []
        for entry in response:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals
            users.append({
                "user": entry['attributes']['uid'][0],
                "display_name": entry['attributes'].get('displayName') or "",
                "locked": entry['attributes'].get('pwdAccountLockedTime', []) != [],
            })
        return users
