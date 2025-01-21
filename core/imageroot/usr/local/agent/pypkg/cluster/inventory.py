#
# Copyright (C) 2025 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import sys
import cluster.modules
import agent

def fact_subscription():
    ret = {"subscription": ""}
    rdb = agent.redis_connect(use_replica=True)
    provider = rdb.hget('cluster/subscription', 'provider')
    if provider == "nsent":
        return "enterprise"
    elif provider == "nscom":
        return "community"
    else:
        return "disabled"


def fact_admins():
    rdb = agent.redis_connect(privileged=True)
    users = {"count": 0, "2fa": 0}
    # get all user keys
    for ukey in rdb.scan_iter('user/*'):
        users["count"] += 1
        two_fa = rdb.hget(ukey, '2fa')
        if two_fa:
            users["2fa"] += 1
    return users


def cluster_facts():
    result = agent.tasks.run(
        agent_id='cluster',
        action='get-facts',
        extra={
            'isNotificationHidden': True,
        },
        endpoint="redis://cluster-leader"
    )

    if result['exit_code'] == 0:
        return result['output']
    else:
        print(agent.SD_ERR + f"cluster/get-facts failed", file=sys.stderr)
        return None
 
def modules_facts():
    ret = []
    rdb = agent.redis_connect(use_replica=True)
    installed_modules = cluster.modules.list_installed(rdb, skip_core_modules = False)
    instance_list = list(mi for module_instances in installed_modules.values() for mi in module_instances)

    # Collect facts from modules:
    for module in instance_list:
        minfo = {"id": module.get('id'), "version": module.get("version"), "name": module.get("module"), "node": module.get("node") } 

        try:
            list_actions_result = agent.tasks.run(
                agent_id='module/' + module.get('id'),
                action='list-actions',
                extra={
                    'isNotificationHidden': True,
                },
                endpoint="redis://cluster-leader",
            )

            if not list_actions_result or list_actions_result['exit_code'] != 0:
                print(agent.SD_WARNING + f"module/{module.get('id')}/list-actions failed", file=sys.stderr)
                raise Exception("list-actions failed")
            
            if 'get-facts' in list_actions_result.get('output', []):

                get_facts_result = agent.tasks.run(
                    agent_id='module/' + module.get('id'),
                    action='get-facts',
                    extra={ 'isNotificationHidden': True },
                    endpoint="redis://cluster-leader",
                )

                if get_facts_result['exit_code'] != 0:
                    print(agent.SD_WARNING, f"get-facts failed for {module.get('id')}", file=sys.stderr)
                    raise Exception("get-facts failed")

                ofacts = get_facts_result['output']
                # merge info from get-facts
                minfo.update(ofacts)
        except:
            pass
        finally:
            ret.append(minfo)

    return ret
