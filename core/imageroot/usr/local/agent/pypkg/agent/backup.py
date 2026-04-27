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

def resolve_gateway_address(rdb, dest_id, protocol='webdav'):
    return [(f"http://127.0.0.1:4694/{dest_id}", os.environ["NODE_ID"])]

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
                hdrs = {"Accept": "application/vnd.x.restic.rest.v2"}
                req = rses.head(rest_url + '/' + probe_path + '/config', headers=hdrs)
                if req.status_code in [404, 200]:
                    return rest_url
                else:
                    req.raise_for_status()
            except requests.RequestException as ex:
                probe_results[node_id] = ex
    raise GatewayProbeException(probe_results)
