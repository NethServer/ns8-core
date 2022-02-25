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

