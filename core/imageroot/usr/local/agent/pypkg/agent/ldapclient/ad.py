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

from datetime import timedelta
import datetime
import ldap3
from .exceptions import LdapclientEntryNotFound
from .base import LdapclientBase

class LdapclientAd(LdapclientBase):
    def _get_dn_attributes(self, dn, lfilter='(objectClass=*)', attributes=[ldap3.ALL_ATTRIBUTES, ldap3.ALL_OPERATIONAL_ATTRIBUTES]):
        response = self.ldapconn.search(dn, lfilter, ldap3.BASE, attributes=attributes)[2]

        for entry in response:
            if entry['type'] == 'searchResEntry':
                return entry['attributes']

        raise LdapclientEntryNotFound()

    def get_max_pwd_age(self):
        response = self.ldapconn.search(self.base_dn, '(objectClass=domainDNS)', attributes=['maxPwdAge'])
        if response[0]:
            result = [entry for entry in response[2] if entry['type'] == 'searchResEntry']
        else:
            return None
        if not result:
            return None
        try:
            max_pwd_age = result[0]['attributes']['maxPwdAge'].total_seconds()
            if max_pwd_age < 86400000000000:
                return max_pwd_age
        except:
            pass
        return None

    def get_group(self, group):
        # Escape group string to build the filter assertion:
        escgroup = ldap3.utils.conv.escape_filter_chars(group)
        response = self.ldapconn.search(self.base_dn, f'(&(objectClass=group)(sAMAccountName={escgroup}){self._get_groups_search_filter_clause()})',
            attributes=['member', 'description', 'sAMAccountName'],
        )[2]

        def lget_user(dn):
            user = self._get_dn_attributes(dn, attributes=['sAMAccountName', 'displayName'])

            return {
                'user': user['sAMAccountName'],
                'display_name': user.get('displayName') or "",
            }

        for entry in response:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals
            return {
                "group": entry['attributes']['sAMAccountName'],
                "description": __class__.get_entry_description(entry),
                "users": [lget_user(dn) for dn in entry['attributes']['member']],
            }

        raise LdapclientEntryNotFound()

    def list_groups(self):
        response = self.ldapconn.search(self.base_dn, f'(&(objectClass=group){self._get_groups_search_filter_clause()})',
            attributes=['cn','member','description', 'sAMAccountName'],
        )[2]

        group_entry_generator = self.ldapconn.extend.standard.paged_search(
            search_base = self.base_dn,
            search_filter = f'(&(objectClass=group){self._get_groups_search_filter_clause()})',
            search_scope = ldap3.SUBTREE,
            attributes = ['cn','member','description', 'sAMAccountName'],
            paged_size = 900,
            generator=True,
        )

        groups = []
        for entry in group_entry_generator:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals
            groups.append({
                "group": entry['attributes']['sAMAccountName'],
                "description": __class__.get_entry_description(entry),
                "users": [dn.split(',', 1)[0].removeprefix('CN=') for dn in entry['attributes']['member']],
            })
        return groups

    def get_user(self, user):
        entry = self.get_user_entry(user)

        def lget_group(dn):
            group = self._get_dn_attributes(dn, attributes=['sAMAccountName', 'description'])

            try:
                description = group['description'][0]
            except IndexError:
                description = ""

            return {
                "group": group['sAMAccountName'],
                "description": description,
            }

        return {
            "user": entry['attributes']['sAMAccountName'],
            "display_name": entry['attributes'].get('displayName') or "",
            "groups": [lget_group(dn) for dn in entry['attributes']['memberOf']],
            "locked": bool(entry['attributes']['userAccountControl'] & 0x2), # ACCOUNTDISABLE
        }

    def get_user_entry(self, user, lextra_attributes=[]):
        # Escape user string to build the filter assertion:
        escuser = ldap3.utils.conv.escape_filter_chars(user)
        response = self.ldapconn.search(self.base_dn, f'(&(objectClass=user)(sAMAccountName={escuser}){self._get_users_search_filter_clause()})',
            attributes=['displayName', 'sAMAccountName', 'memberOf', 'userAccountControl'] + self.filter_schema_attributes(lextra_attributes),
        )[2]

        for entry in response:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals

            return entry

        raise LdapclientEntryNotFound()

    def list_users(self, extra_info=False):
        attributes = ['displayName', 'sAMAccountName', 'userAccountControl']
        if extra_info:
            attributes += ['whenCreated', 'pwdLastSet', 'mail']
        user_entry_generator = self.ldapconn.extend.standard.paged_search(
            search_base = self.base_dn,
            search_filter = f'(&(objectClass=user)(objectCategory=person){self._get_users_search_filter_clause()})',
            search_scope = ldap3.SUBTREE,
            attributes = attributes,
            paged_size = 900,
            generator=True,
        )

        users = []
        max_pwd_age = self.get_max_pwd_age()
        today = datetime.datetime.now(datetime.timezone.utc)
        for entry in user_entry_generator:
            if entry['type'] != 'searchResEntry':
                continue # ignore referrals
            user = {
                "user": entry['attributes']['sAMAccountName'],
                "display_name": entry['attributes'].get('displayName') or "",
                "locked": bool(entry['attributes']['userAccountControl'] & 0x2), # ACCOUNTDISABLE
            }
            if extra_info:
                pwd_changed_time = entry['attributes'].get('pwdLastSet', entry['attributes'].get('whenCreated', None))
                expire = (entry['attributes']['userAccountControl'] & 0x10000 == 0) # DONT_EXPIRE_PASSWORD
                if expire and max_pwd_age:
                    expiry_date = pwd_changed_time + timedelta(seconds=max_pwd_age)
                    user['password_expiration'] = int(expiry_date.timestamp())
                    user['expired'] = today > expiry_date
                else:
                    user['password_expiration'] = -1
                    user['expired'] = False
                # mail can be a string or an empty array, just treat ans empty arrays as an empty string
                user["mail"] = entry['attributes'].get('mail') if entry['attributes'].get('mail') else ""
                # In AD, pwdLastSet = 0 means that the user must change the password at next logon
                # The timestamp 0 in Windows epoch is -11644473600 in Unix epoch
                user['must_change'] = (pwd_changed_time.timestamp() == -11644473600)
            users.append(user)

        return users
