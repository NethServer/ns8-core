#
# Copyright (C) 2025 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

import sys
import subprocess
import json
import os
from configparser import ConfigParser

def _parse_config() -> ConfigParser:
    config = ConfigParser(
        delimiters=("="),
        comment_prefixes=("#"),
        inline_comment_prefixes=("#"),
        empty_lines_in_values=False,
    )
    try:
        config.read("/etc/nethserver/volumes.conf")
    except FileNotFoundError:
        pass
    return config

def get_application_types() -> list:
    oconfig = _parse_config()
    apptypes = set()
    for section in oconfig:
        app_type = section.rstrip('0123456789')
        apptypes.add(app_type)
    apptypes.discard('DEFAULT')
    return sorted(apptypes)

def get_configuration(module_id:str) -> dict:
    """Return the volume configuration for the given module_id by reading
    the config file /etc/nethserver/volumes.conf."""
    app_type = module_id.rstrip('0123456789')
    oconfig = _parse_config()
    if module_id in oconfig:
        section_key = module_id
    elif app_type in oconfig:
        section_key = app_type
    else:
        return {} # config not found
    # Check if the paths are valid mountpoints
    valid_paths = [bp["path"] for bp in get_base_paths()]
    volmap = {}
    for volname, volpath in oconfig[section_key].items():
        volpath = volpath.rstrip("/")
        if volpath in valid_paths:
            volmap[volname] = volpath
        else:
            print(f"volumes.py: invalid base path {volpath} (volume '{volname}', app '{section_key}')", file=sys.stderr)
    return volmap

def get_base_paths() -> list:
    home_basedir = os.getenv("HOME_BASEDIR", "/home")
    with subprocess.Popen(["findmnt", "--real", "--json", "--bytes", "--output", "SOURCE,TARGET,SIZE,USED,LABEL", "--nofsroot", "--type", "nofat,novfat"], stdout=subprocess.PIPE) as hproc:
        ofindmnt = json.load(hproc.stdout)
    orootfs = ofindmnt["filesystems"][0]
    paths = []
    for ofs in orootfs.get("children", []):
        if ofs["source"] == orootfs["source"]:
            continue # ignore bind mounts
        elif ofs["target"].startswith("/boot"):
            continue # ignore /boot* mountpoints
        elif ofs["target"] in ["/", home_basedir]:
            continue # ignore HOME_BASEDIR and root volume
        opath = {
            "path": ofs["target"],
            "label": ofs["label"] or "",
        }
        if "size" in ofs:
            opath["size"] = ofs["size"]
            opath["used"] = ofs["used"]
            opath["available"] = ofs["size"] - ofs["used"]
        paths.append(opath)
    return paths

if __name__ == "__main__":
    json.dump(get_base_paths(), fp=sys.stdout)

