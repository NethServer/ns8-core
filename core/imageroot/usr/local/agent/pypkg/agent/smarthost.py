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

class SmartHost:

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
                "enabled": conf['enabled'] == "1",
                "tls": conf['tls'] == "1",
                "tls_verify": conf['tls_verify'] == "1"
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
   json.dumps(sp.load_smtp())
