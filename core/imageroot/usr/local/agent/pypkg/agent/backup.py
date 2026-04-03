#
# Copyright (C) 2026 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import time
import os
import json
import requests
from .safeio import exclusive_file_lock

CACHE_PATH = './.gateway_resolver.cache'
CACHE_EXPIRY = 3600 # (seconds) per-entry

def update_gateway_resolver_cache(dest_id, node_id):
    """Store the preferred node_id for dest_id in a local directory file
    cache."""
    tsnow = int(time.time())
    with exclusive_file_lock(CACHE_PATH):
        fo = None
        try:
            fo = open(CACHE_PATH, "r")
            ocache = json.load(fo)
            tsold = ocache[dest_id][1]
        except (FileNotFoundError, json.JSONDecodeError):
            ocache = {} # Init cache
            tsold = 0
        except KeyError:
            tsold = 0 # Cache miss
        finally:
            if fo is not None:
                fo.close()
        needs_update = False
        for key, entry in list(ocache.items()):
            tsold = entry[1]
            # Expiry check for each entry
            if tsold + (CACHE_EXPIRY) < tsnow:
                needs_update = True
                if key == dest_id:
                    ocache[dest_id] = [str(node_id), tsnow]
                else:
                    # Drop expired entry
                    del ocache[key]
        if dest_id not in ocache or ocache[dest_id][0] != str(node_id):
            ocache[dest_id] = [str(node_id), tsnow]
            needs_update = True
        if needs_update:
            # Rewrite cache file if something was changed:
            with open(CACHE_PATH, 'w') as fo:
                json.dump(ocache, fp=fo, separators=(",", ":"))

def resolve_gateway_address(rdb, dest_id, protocol='webdav'):
    """Return a list of rclone-gateway URLs for WebDAV port, high priority
    first established by a cache file in the current working directory."""

    if protocol not in ["webdav", "rest"]:
        raise Exception("Invalid protocol: " + str(protocol))

    cache_maxage = 3600 # 1 hour
    tsnow = int(time.time())

    # Unordered set of node IDs
    node_ids = set(rdb.hvals("cluster/module_node"))

    # Output: ordered list of node IDs
    prio_ids = []

    # First choice, read node ID from cache file:
    try:
        last_id = None
        tsold = None
        fo = open(CACHE_PATH, 'r')
        ocache = json.load(fo)
        last_id, tsold = ocache[dest_id]
        if tsold + (CACHE_EXPIRY) < tsnow:
            # Cache entry is valid
            node_ids.discard(last_id)
            prio_ids.append(last_id)
    except (FileNotFoundError, KeyError, json.JSONDecodeError):
        pass

    # Second choice, the local node:
    local_id = os.environ["NODE_ID"]
    node_ids.discard(local_id)
    if local_id not in prio_ids:
        prio_ids.append(local_id)

    # Third choice, the leader node, VPN hub:
    leader_id = rdb.hget("cluster/environment", "NODE_ID")
    node_ids.discard(leader_id)
    if leader_id not in prio_ids:
        prio_ids.append(leader_id)

    # Fourth choice, append remaining nodes
    prio_ids += list(node_ids)

    gwlist = []
    for node_id in prio_ids:
        if node_id == local_id:
            node_vpn_ip = '127.0.0.1'
        else:
            node_vpn_ip = rdb.hget(f"node/{node_id}/vpn", "ip_address")

        if node_vpn_ip:
            if protocol == 'webdav':
                gwlist.append((f"http://{node_vpn_ip}:4694/{dest_id}", node_id))
            elif protocol == 'rest':
                gwlist.append((f"http://{node_vpn_ip}:4695/{dest_id}", node_id))

    return gwlist

class TimeoutSession(requests.Session):
    def __init__(self, timeout=10, auth=None):
        super().__init__()
        self.timeout = timeout  # (connect, read) or single float
        if auth == None:
            self.auth = (os.environ["REDIS_USER"], os.environ["REDIS_PASSWORD"])
        else:
            self.auth = auth
    def request(self, method, url, **kwargs):
        kwargs.setdefault("timeout", self.timeout)
        return super().request(method, url, **kwargs)

class GatewayProbeException(Exception):
    def __init__(self, probes: dict):
        super().__init__()
        self.probes = probes

    def __str__(self):
        parts = []
        for node_id, exc in self.probes.items():
            reason = self._format_exception(exc)
            parts.append(f"node {node_id}: {reason}")
        return "Gateway probe failed: " + "; ".join(parts)

    @staticmethod
    def _format_exception(exc):
        # requests exceptions may have response/status info
        if hasattr(exc, "response") and exc.response is not None:
            req = exc.response.request
            if req is not None:
                return f"{req.url} -> HTTP {exc.response.status_code} - {exc.response.reason}"
            else:
                return f"HTTP {exc.response.status_code} - {exc.response.reason}"
        return str(exc)

def probe_restic_gateway(rdb, dest_id, probe_path, http_auth=None):
    """Return a working restic REST backend URL by trying different
    rclone-gateway instances running on each cluster node."""
    probe_results = {}
    with TimeoutSession(timeout=(3, 10), auth=http_auth) as rses:
        for rest_url, node_id in resolve_gateway_address(rdb, dest_id, protocol='rest'):
            try:
                req = rses.head(rest_url + '/' + probe_path + '/config')
                if req.status_code in [404, 200]:
                    update_gateway_resolver_cache(dest_id, node_id)
                    return rest_url
                else:
                    req.raise_for_status()
            except requests.RequestException as ex:
                probe_results[node_id] = ex
    raise GatewayProbeException(probe_results)
