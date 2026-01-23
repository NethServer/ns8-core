#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

# Fact collection helpers

import hashlib
import ipaddress

PSEUDO_ENFORCING = True
PSEUDO_SEED = '0000'

def init_pseudonymization(enforce, rdb):
    global PSEUDO_ENFORCING, PSEUDO_SEED
    seed = rdb.get('cluster/anon_seed')
    if not seed:
        raise Exception("The cluster/anon_seed value is not set")
    PSEUDO_ENFORCING = enforce
    PSEUDO_SEED = seed

def has_subscription(rdb):
    provider = rdb.hget('cluster/subscription', 'provider')
    return provider in ["nscom", "nsent"]

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

def pseudo_ip(val):
    """Calculate a stable pseudonym of the given IPv4 or IPv6 address,
    preserving only private vs public scope
    """
    if not val or not PSEUDO_ENFORCING:
        return val

    try:
        ip = ipaddress.ip_address(val)
    except ValueError:
        return val

    digest = hashlib.md5((PSEUDO_SEED + ip.exploded).encode('utf-8')).digest()

    if isinstance(ip, ipaddress.IPv4Address):
        if ip.is_private:
            # 10.0.0.0/8
            host = int.from_bytes(digest[:3], byteorder='big')
            pseudo_int = (10 << 24) | host
        else:
            # 1.0.0.0/8 (public)
            host = int.from_bytes(digest[:3], byteorder='big')
            pseudo_int = (1 << 24) | host

        return str(ipaddress.IPv4Address(pseudo_int))

    else:
        if ip.is_private:
            # fc00::/7 (ULA)
            host = int.from_bytes(digest[:15], byteorder='big')
            pseudo_int = (0xfc << 120) | host
        else:
            # 2000::/3 (global unicast)
            host = int.from_bytes(digest[:15], byteorder='big')
            pseudo_int = (0x2 << 124) | host

        return str(ipaddress.IPv6Address(pseudo_int))
