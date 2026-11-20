#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import os
import json
import agent

def is_mail_profile(smarthost):
    """Return True if the smarthost settings are the ones the UI writes
    when a Mail module of the cluster is chosen"""
    return smarthost.get('enabled') == '1' \
        and str(smarthost.get('port')) == '25' \
        and not smarthost.get('username') \
        and smarthost.get('encrypt_smtp', 'none') == 'none'

def find_mail_uuid(rdb, smarthost):
    """Return the UUID of the Mail module the smarthost points to, or
    None if it is not a Mail module of the cluster"""
    if not is_mail_profile(smarthost):
        return None
    for xprov in agent.list_service_providers(rdb, 'submission', 'tcp'):
        if xprov.get('host') == smarthost.get('host'):
            return xprov.get('module_uuid')
    return None

def sync_mail_host(rdb):
    """Point the smarthost to the current host of the linked Mail module.
    Return the new host, or None if nothing changed"""
    # The cluster agent of every node gets the event, but only the leader Redis is writable
    if os.environ.get('NODE_ID') != rdb.hget('cluster/environment', 'NODE_ID'):
        return None
    smarthost = rdb.hgetall('cluster/smarthost')
    mail_uuid = smarthost.get('mail_uuid')
    if not mail_uuid or not is_mail_profile(smarthost):
        return None
    providers = agent.list_service_providers(rdb, 'submission', 'tcp', filters={'module_uuid': mail_uuid})
    host = providers[0].get('host') if providers else None
    if not host or host == smarthost.get('host'):
        return None
    trx = rdb.pipeline()
    trx.hset('cluster/smarthost', 'host', host)
    trx.publish('cluster/event/smarthost-changed', json.dumps({
        'key': 'cluster/smarthost',
    }))
    trx.execute()
    return host
