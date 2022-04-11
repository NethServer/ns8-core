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

import redis
import os

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
                'base_dn',
                'bind_dn',
                'bind_password',
            ]

            for key in keys:
                conf[key] = self.domains[domain][key]

            conf['port'] = self.domains[domain]['listen_port']

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

        domains = {}
        for ksrv in rdb.scan_iter('*/srv/tcp/ldap'):
            dhx = rdb.hgetall(ksrv)
            if 'domain' in dhx:
                domains.setdefault(dhx['domain'], dhx)

            if dhx['domain'] in domain_port:
                domains[dhx['domain']].setdefault('listen_port', domain_port[dhx['domain']])

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

if __name__ == '__main__':
    lp = Ldapproxy()
    for domain in lp.get_domains_list():
        print(domain, f'node/{lp.node_id}', lp.get_domain(domain))
