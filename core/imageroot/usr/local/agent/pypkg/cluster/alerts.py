#
# Copyright (C) 2025 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import requests

def _decorate_with_our_attributes(oalert):
    if "alertname" in oalert.get("labels", {}):
        oalert["name"] = oalert["labels"]["alertname"]
    else:
        return # ignore alert
    target_type = oalert.get("labels", {}).get("target_type")
    if oalert["name"] == "backup_failed":
        # A backup has many applications from different nodes, so it
        # defines a category on its own:
        oalert["category"] = "backup"
    elif target_type in ["traefik", "node"]:
        # If target is node_exporter or traefik, assign "node" category:
        oalert["category"] = "node"
        oalert["node_id"] = int(oalert["labels"].get("node", "0"))
    elif 'module_id' in oalert.get("labels", {}):
        oalert["category"] = "application" # Fallback category
        oalert["module_id"] = oalert["labels"]["module_id"]

def fetch_alerts():
    """Fetch active alerts from local Alertmanager."""
    alertmanager_url = "http://127.0.0.1:9093/api/v2/alerts"
    r = requests.get(alertmanager_url)
    r.raise_for_status()
    result = []
    for oalert in r.json():
        _decorate_with_our_attributes(oalert)
        if not "name" in oalert:
            continue
        elif not "category" in oalert:
            continue
        result.append(oalert)
    return result
