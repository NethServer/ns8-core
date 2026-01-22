#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

# Fact collection helpers

import hashlib

PSEUDO_ENFORCING = True
PSEUDO_SEED = '0000'

def init_pseudonymization(enforce, rdb):
    global PSEUDO_ENFORCING, PSEUDO_SEED
    seed = rdb.get('cluster/anon_seed')
    if not seed:
        raise Exception("The cluster/anon_seed value is not set")
    PSEUDO_ENFORCING = enforce
    PSEUDO_SEED = seed

def pseudo_string(val, maxlen=12):
    """Calculate a stable pseudonym of the given string"""
    if val and PSEUDO_ENFORCING:
        hashed_val = hashlib.md5((PSEUDO_SEED + val).encode('utf-8')).hexdigest()
        return hashed_val[0:maxlen]
    else:
        return val

def pseudo_domain(val):
    """Calculate a stable pseudonym of the given domain, keeping the TLD in clear text"""
    try:
        domain, suffix = val.rsplit(".", 1)
        return pseudo_string(domain, 8) + '.' + suffix
    except ValueError:
        return pseudo_string(val)
