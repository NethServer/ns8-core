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

import fnmatch
import agent

def _lookup_grant_on(rdb, agent_pattern):
    """
    Find matching agents by pattern
    """
    results = []
    for key in list(rdb.scan_iter(f'{agent_pattern}/roles/*')):
        pos = key.find('/roles/')
        if pos > -1:
            results.append(key[:pos])
    return results

def _lookup_grant_to(rdb, role_pattern):
    """
    Find matching roles by pattern
    """
    results = []
    for key in list(rdb.scan_iter(f'*/roles/{role_pattern}')):
        pos = key.find('/roles/')
        if pos > -1:
            results.append(key[pos+7:])
    return results

def _change_role_definition(rdb, revoke, action_clause, on_clause, to_clause):
    if "*" in on_clause:
        on_clause_list = _lookup_grant_on(rdb, on_clause)
    else:
        on_clause_list = [on_clause]

    if "*" in to_clause:
        to_clause_list = _lookup_grant_to(rdb, to_clause)
    else:
        to_clause_list = [to_clause]

    pipe = rdb.pipeline()
    for onk in on_clause_list:
        for tok in to_clause_list:
            if revoke:
                pipe.srem(f'{onk}/roles/{tok}', action_clause)
            else:
                pipe.sadd(f'{onk}/roles/{tok}', action_clause)
    pipe.execute()

def grant(rdb, action_clause, on_clause, to_clause):
    """
    Grant ACTION ON agent TO role.

    If wildcard "*" is present in the ACTION, ON and TO clauses
    lookup matching keys in Redis DB.
    """
    return _change_role_definition(rdb, False, action_clause, on_clause, to_clause)

def revoke(rdb, action_clause, on_clause, to_clause):
    """
    Revoke ACTION ON agent TO role.

    See grant_action()
    """
    return _change_role_definition(rdb, True, action_clause, on_clause, to_clause)

def alter_user(rdb, user, revoke, role, on_clause):
    """
    Grant/Revoke the ROLE ON some module to USER

    The module can be cluster, node/X, module/modY etc.
    The ROLE must be already defined in the given module.
    """

    pipe = rdb.pipeline()
    for key in rdb.scan_iter(f'{on_clause}/roles/{role}'):
        pos = key.find(f'/roles/{role}')
        agent_id = key[:pos]
        if revoke:
            pipe.hdel(f'roles/{user}', agent_id, role)
        else:
            pipe.hset(f'roles/{user}', agent_id, role)

    pipe.execute()

def refresh_permissions(rdb):
    """Scan the stored authorizations and re-apply permissions to all modules"""
    for kauth_module in rdb.scan_iter('cluster/authorizations/module/*'):
        kmid = kauth_module.removeprefix('cluster/authorizations/module/')
        authorizations = rdb.smembers(kauth_module) or []
        mnode_id = rdb.hget('cluster/module_node', kmid)
        add_module_permissions(rdb, kmid, authorizations, mnode_id)

def add_module_permissions(rdb, module_id, authorizations, node_id):
    """Parse authorizations and grant permissions to module_id"""
    for authz in authorizations:
        xagent, xrole = authz.split(':') # e.g. traefik@node:routeadm

        if xagent.endswith("@any"):
            agent_selector = 'module/' + xagent.removesuffix("@any") + '*' # wildcard allowed in on_clause
        elif xagent == 'self':
            agent_selector = 'module/' + module_id
        else:
            agent_selector = agent.resolve_agent_id(xagent, node_id=node_id)

        alter_user(rdb,
            user=f'module/{module_id}',
            revoke=False,
            role=xrole,
            on_clause=agent_selector,
        )

def check_authorizations_sanity(authorizations):
    for authz in authorizations:
        # Validation assertions:
        contains_forbidden_char = ('*' in authz or '[' in authz or '?' in authz)
        lacks_mandatory_separator = (not ':' in authz)
        grants_too_much_permissions = (authz == 'cluster:owner')
        if contains_forbidden_char or lacks_mandatory_separator or grants_too_much_permissions:
            return False

    return True

def save_acls(rdb):
    """
    Save ACLs on the cluster leader disk
    Copy current ACLs to cluster/acls key
    Publish acl-changed event: worker nodes will run acl-load and save ACLs on their disks too

    This function can be executed only on the leader node
    """

    to_skip = ["default", "cluster", "api-server"]

    acl_list = rdb.acl_list()

    trx = rdb.pipeline()

    # Cleanup ACLs
    trx.delete("cluster/acls")

    # Skip ACLs which should be different on each node
    for acl in acl_list:
        user = acl.split(" ",2)[1]
        if user in to_skip:
            continue
        trx.sadd("cluster/acls", acl)

    trx.acl_save()
    trx.publish('cluster/event/acl-changed', 'cluster/acls')
    trx.execute()
