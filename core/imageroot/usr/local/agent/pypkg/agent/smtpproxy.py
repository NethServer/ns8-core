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
import json

class smtpproxy:

    def load_smtp(self):
        rdb = self._redis_connect()
        conf = rdb.hgetall('cluster/smarthost')

        if not conf:
            data={
                "port": 587,
                "host": "",
                "username":"",
                "password": "",
                "enabled": False,
                "tls": True,
                "tls_verify": True
            }
        else:
            data={
                "port": int(conf['port']),
                "host": conf['host'],
                "username": conf['username'],
                "password": conf['password'],
                "enabled": conf['enabled'] == "on",
                "tls": conf['tls'] == "on",
                "tls_verify": conf['tls_verify'] == "on"
            }
        return data
        
    def _redis_connect(self):
        return redis.Redis(
            host='127.0.0.1',
            port=6379,
            username='default',
            password='default',
            decode_responses=True,
        )

if __name__ == '__main__':
    sp = smtpproxy()
    print(json.dumps(sp.load_smtp()))
