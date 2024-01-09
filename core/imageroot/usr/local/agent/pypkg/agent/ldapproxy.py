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

import sys
import agent
import redis
import os
import cluster.userdomains

class Ldapproxy:

    def __init__(self, node_id=None):
        self.domains = None
        if node_id is None:
            node_id = os.environ['NODE_ID']
        self.node_id = int(node_id)

    def get_domains_list(self):
        if self.domains is None:
            self.domains = self._load_domains()
        return list(self.domains.keys())

    def get_domain(self, domain):
        if self.domains is None:
            self.domains = self._load_domains()
        conf = {
            'host': '127.0.0.1',
        }
        try:
            keys = [
                'schema',
                'location',
                'base_dn',
                'bind_dn',
                'bind_password',
            ]

            for key in keys:
                conf[key] = self.domains[domain][key]

            conf['port'] = self.domains[domain]['listen_port']
            if conf['schema'] == 'ad':
                conf['hidden_users'] = [
                    'Guest',
                    'krbtgt',
                    'ldapservice',
                ]
                conf['hidden_groups'] = [
                    'DnsUpdateProxy',
                    'Domain Computers',
                    'Domain Controllers',
                    'Domain Guests',
                    'Domain Users',
                    'Group Policy Creator Owners',
                    'Read-only Domain Controllers',
                    'Protected Users',
                ]
            elif conf['schema'] == 'rfc2307':
                conf['hidden_users'] = []
                conf['hidden_groups'] = ['locals']
            else:
                conf['hidden_users'] = []
                conf['hidden_groups'] = []

        except:
            return None

        return conf

    def _load_domains(self):
        rdb = self._redis_connect()

        ldapproxy_instance = rdb.get(f'node/{self.node_id}/default_instance/ldapproxy')
        assert(not ldapproxy_instance is None)

        domain_port = {}
        try:
            domain_port.update(rdb.hgetall(f'module/{ldapproxy_instance}/data/domain_port'))
        except TypeError:
            pass

        # Retrieve the list of bound user domains, to check warning
        # conditions:
        bound_domain_list = agent.get_bound_domain_list(rdb)

        domains = {}
        configured_domains = cluster.userdomains.list_domains(rdb)
        for domain in configured_domains:
            if "MODULE_ID" in os.environ and domain not in bound_domain_list:
                # Warn only if the module is running under a module environment
                print(agent.SD_WARNING + \
                    f'agent.ldapproxy: domain {domain} should not be used by ' + \
                    f'{os.getenv("MODULE_ID")}. Invoke agent.bind_user_domains(' + \
                    f'["{domain}"]) to fix this warning.', file=sys.stderr)
            dhx = configured_domains[domain]
            domains.setdefault(domain, dhx)

            if domain in domain_port:
                domains[domain].setdefault('listen_port', domain_port[domain])

        rdb.close()
        return domains

    def _redis_connect(self):
        return redis.Redis(
            host='127.0.0.1',
            port=6379,
            username='default',
            password='default',
            decode_responses=True,
        )

    def get_ldap_users_search_filter_clause(self, domain):
        mdom = self.get_domain(domain)
        if not mdom['hidden_users']:
            return ""

        filter_clause = ""
        if mdom['schema'] == 'ad':
            uattr = 'sAMAccountName'
        elif mdom['schema'] == 'rfc2307':
            uattr = 'uid'
        else:
            return ""

        for user in mdom['hidden_users']:
            filter_clause += f"({uattr}={user})"

        return f"(!(|{filter_clause}))"


    def get_ldap_groups_search_filter_clause(self, domain):
        mdom = self.get_domain(domain)

        filter_clause = ""
        if mdom['schema'] == 'ad':
            uattr = 'sAMAccountName'
            # filter out non-global groups. See
            # https://social.technet.microsoft.com/wiki/contents/articles/5392.active-directory-ldap-syntax-filters.aspx
            filter_clause += "(!(groupType:1.2.840.113556.1.4.803:=2))"
        elif mdom['schema'] == 'rfc2307':
            uattr = 'cn'
        else:
            return ""

        for user in mdom['hidden_groups']:
            filter_clause += f"({uattr}={user})"

        if not filter_clause:
            return "" # do not wrap an empty clause!

        return f"(!(|{filter_clause}))"

if __name__ == '__main__':
    lp = Ldapproxy()
    for domain in lp.get_domains_list():
        print(domain, f'node/{lp.node_id}', lp.get_domain(domain))
