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

class LdapclientBase:

    def __init__(self, host, port, bind_dn, bind_password, base_dn, schema, hidden_users, hidden_groups, **kwargs):
        self.host = host
        self.port = int(port)
        self.bind_dn = bind_dn
        self.bind_password = bind_password
        self.base_dn = base_dn
        self.schema = schema
        self.hidden_users = hidden_users
        self.hidden_groups = hidden_groups
        self.extra_args = kwargs

        self.ldapsrv = ldap3.Server(self.host, port=self.port)
        self.ldapconn = ldap3.Connection(self.ldapsrv,
            user=self.bind_dn or None,
            password=self.bind_password or None,
            client_strategy=ldap3.SAFE_SYNC,
            auto_bind=True,
            raise_exceptions=True,
            auto_range=True,
            auto_referrals=False,
        )

    def __enter__(self):
        pass

    def __exit__(self, exc_type, exc_value, traceback):
        if not self.ldapconn.closed:
            self.ldapconn.unbind()

    def get_entry_description(entry):
        try:
            return entry['attributes']['description'][0]
        except IndexError:
            return ""

    def _get_groups_search_filter_clause(self):
        filter_clause = ""
        if self.schema == 'ad':
            uattr = 'sAMAccountName'
            # filter out non-global groups. See
            # https://social.technet.microsoft.com/wiki/contents/articles/5392.active-directory-ldap-syntax-filters.aspx
            filter_clause += "(!(groupType:1.2.840.113556.1.4.803:=2))"
        elif self.schema == 'rfc2307':
            uattr = 'cn'
        else:
            return ""
        for user in self.hidden_groups:
            filter_clause += f"({uattr}={user})"

        if not filter_clause:
            return "" # do not wrap an empty clause!

        return f"(!(|{filter_clause}))"


    def _get_users_search_filter_clause(self):
        if not self.hidden_users:
            return ""

        filter_clause = ""
        if self.schema == 'ad':
            uattr = 'sAMAccountName'
        elif self.schema == 'rfc2307':
            uattr = 'uid'
        else:
            return ""
        for user in self.hidden_users:
            filter_clause += f"({uattr}={user})"

        return f"(!(|{filter_clause}))"

    def filter_schema_attributes(self, attr_list):
        """Reduce the given attr_list by filtering out attribute names not
           declared in the server schema."""
        return [aname for aname in attr_list if aname in self.ldapsrv.schema.attribute_types]
