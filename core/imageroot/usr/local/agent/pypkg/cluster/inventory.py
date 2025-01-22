#
# Copyright (C) 2025 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import os
import agent

def fact_subscription(rdb):
    # Get subscription status
    provider = rdb.hget('cluster/subscription', 'provider')
    if provider == "nsent":
        return "enterprise"
    elif provider == "nscom":
        return "community"
    else:
        return "disabled"

def fact_admins(rdb):
    # Count admin users and users with 2FA enabled
    users = {"count": 0, "2fa": 0}
    for ukey in rdb.scan_iter('user/*'):
        users["count"] += 1
        two_fa = rdb.hget(ukey, '2fa')
        if two_fa:
            users["2fa"] += 1
    return users

def fact_repositories(rdb):
    # List enabled repositories
    ret = []
    for m in rdb.scan_iter('cluster/repository/*'):
        repo = rdb.hgetall(m)
        if repo.get("status") == "1":
            ret.append(os.path.basename(m))
    return ret

def fact_smarthost(rdb):
    # Get smarthost configuration
    manual_configuration = True
    smtp = agent.get_smarthost_settings(rdb)
    # we query all the list of service provider
    for key in agent.list_service_providers(rdb,'submission','tcp'):
        if manual_configuration == True:
            # no need to check the rest of the list if manual_configuration == False
            manual_configuration = False if key['host'] == smtp['host'] and str(smtp['port']) == "25" else True
    return{"enabled": smtp["enabled"], "manual_configuration": manual_configuration}