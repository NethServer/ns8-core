#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import json
import pwd
import sys
import agent

def basic_auth():
    """Return the (username, password) of the local Prometheus and
    Alertmanager, read from the metrics module running on this (leader)
    node, or None if it cannot be read."""
    try:
        metrics_id = agent.redis_connect().get('cluster/default_instance/metrics')
        home = pwd.getpwnam(metrics_id).pw_dir
        with open(f"{home}/.config/state/basic_auth.json") as fp:
            creds = json.load(fp)
        return (creds['username'], creds['password'])
    except Exception as ex:
        print("cluster.metrics.basic_auth() cannot read the metrics password. Reason: ", ex, file=sys.stderr)
        return None
