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
import sys
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

def _lookup_role_targets(rdb, on_clause, role):
    """Return the agent IDs matching ON_CLAUSE that define at least one
    of the comma-separated role names of ROLE"""
    targets = set()
    for xrole in role.split(","):
        # Make sure there is no white space around the role name:
        xrole = xrole.strip()
        if not xrole:
            continue
        for key in rdb.scan_iter(f'{on_clause}/roles/{xrole}'):
            targets.add(key.removesuffix(f'/roles/{xrole}'))
    return targets

def alter_user(rdb, user, revoke, role, on_clause):
    """
    Grant/Revoke the ROLE ON some module to USER

    The module matcher, specified by ON_CLAUSE, can be cluster, node/X,
    module/modY etc. Wildcard match is allowed.

    ROLE is a comma-separated list of role names without spaces. Each role
    name must be already defined in the matching module(s).
    """

    targets = _lookup_role_targets(rdb, on_clause, role)

    pipe = rdb.pipeline()
    for agent_id in targets:
        if revoke:
            pipe.hdel(f'roles/{user}', agent_id)
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

def _resolve_module_permissions(rdb, module_id, authorizations, node_id):
    """Parse authorizations and return the {agent_id: role} map to grant
    to module_id. Redis is only read."""
    permissions = {}
    for authz in authorizations:
        # authz examples: "traefik@node:routeadm", "node:fwadm,portsadm"
        xagent, xrole = authz.split(':')

        if xagent.endswith("@any"):
            agent_selector = 'module/' + xagent.removesuffix("@any") + '*' # wildcard allowed in on_clause
        elif xagent == 'self':
            agent_selector = 'module/' + module_id
        else:
            agent_selector = agent.resolve_agent_id(xagent, node_id=node_id)

        if agent_selector is None:
            print(agent.SD_WARNING + f"Cannot resolve {xagent} for module/{module_id}, authorization {authz} skipped", file=sys.stderr)
            continue

        # If the agent_selector resolves multiple times to the same value,
        # the last authz wins: roles are not accumulated or merged. Still
        # it is allowed to pass a comma-separated list of roles as xrole
        # value.
        for agent_id in _lookup_role_targets(rdb, agent_selector, xrole):
            permissions[agent_id] = xrole
    return permissions

def add_module_permissions(rdb, module_id, authorizations, node_id):
    """Parse authorizations and grant permissions to module_id"""
    permissions = _resolve_module_permissions(rdb, module_id, authorizations, node_id)
    if permissions:
        rdb.hset(f'roles/module/{module_id}', mapping=permissions)

def update_module_permissions(rdb, module_id, old_authorizations, new_authorizations, node_id):
    """Grant the new authorizations to module_id and revoke the permissions
    granted by old authorizations that are no longer requested.
    Permissions granted by other means are kept."""
    old_permissions = _resolve_module_permissions(rdb, module_id, old_authorizations, node_id)
    new_permissions = _resolve_module_permissions(rdb, module_id, new_authorizations, node_id)
    current_permissions = rdb.hgetall(f'roles/module/{module_id}')
    # A field still holding the role of the old label was granted by it.
    # If it holds something else, it was changed by other means: keep it.
    obsolete = [agent_id for agent_id, role in old_permissions.items()
        if agent_id not in new_permissions and current_permissions.get(agent_id) == role]
    trx = rdb.pipeline()
    if obsolete:
        trx.hdel(f'roles/module/{module_id}', *obsolete)
    if new_permissions:
        trx.hset(f'roles/module/{module_id}', mapping=new_permissions)
    trx.execute()

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
