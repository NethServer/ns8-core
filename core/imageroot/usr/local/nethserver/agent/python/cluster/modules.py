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

import os
import json
import agent
import semver
import urllib
import os.path
import urllib.request

def __urljoin(base_path, *args):
    '''replace urllib.parse.joinurl because it doesn't handle multiple parameters
    '''
    parts = [base_path]
    for arg in args:
        # also make sure to escape special chars using quote
        parts.append(urllib.parse.quote(arg))
    return "/".join(part.strip("/") for part in parts)

def __parse_repository_metadata(repository_name, repository_url, repository_updated, repodata):
    modules = []

    try:
        repodata = json.loads(repodata)
    except:
        return modules

    for package in repodata:
        package["repository"] = repository_name
        package["repository_updated"] = repository_updated

        # Set absolute path for logo
        if package["logo"]:
           package["logo"] = __urljoin(repository_url, package["id"], package["logo"])

        # Set absolute path for screenshots
        screenshots = []
        for s in package["screenshots"]:
           screenshots.append(__urljoin(repository_url, package["id"], s))
        package["screenshots"] = screenshots

        modules.append(package)

    return modules


def __list_repository_modules(rdb, repository_name, repository_url):
    key = f'cluster/repository_cache/{repository_name}'
    cache = rdb.hgetall(key)
    if cache:
        return __parse_repository_metadata(repository_name, repository_url, cache["updated"], cache["data"])

    try:
        url = urllib.parse.urljoin(repository_url, "repodata.json")
        with urllib.request.urlopen(url) as req:
            repodata = req.read().decode()
    except:
        # If repository is not accessible or invalid, just return an empty array
        return []

    updated = req.headers.get('Last-Modified', "")
    modules = __parse_repository_metadata(repository_name, repository_url, updated, repodata)
    # Save inside the cache if data is valid
    if modules:
        # Save also repodata file date
        rdb.hset(key, mapping={"data": repodata, "updated": updated})
        # Set cache expiration to 3600 seconds
        rdb.expire(key, 3600)

    return modules

def list_available(rdb):
    modules = []
    # List all modules from enabled repositories
    for m in rdb.scan_iter('cluster/repository/*'):
        repo = rdb.hgetall(m)
        # Skip disabled repositories
        if int(repo["status"]) == 0:
            continue
        
        modules.extend(__list_repository_modules(rdb, os.path.basename(m), repo["url"]))

    return modules

def list_installed(rdb):
    installed = {}

    # Search for installed modules
    for m in rdb.scan_iter('module/*/environment'):
        vars = rdb.hgetall(m)
        url, sep, tag = vars['IMAGE_URL'].partition(":")
        if url not in installed.keys():
            installed[url] = []
        installed[url].append({ 'id': vars["MODULE_ID"], 'node': vars['NODE_ID'], 'digest': vars["IMAGE_DIGEST"], 'source': url, 'version': tag })

    return installed

def list_updates(rdb):
    updates = []
    installed = list_installed(rdb)
    available = list_available(rdb)

    for module in available:
        if module["source"] not in installed.keys():
            continue

        newest_version = None
        for version in module["versions"]:
            v = semver.VersionInfo.parse(version["tag"])
            # Skip testing versions if testing is disabled
            testing = rdb.hget(f'cluster/repository/{module["repository"]}', 'testing')
            if int(testing) == 0 and not v.prerelease is None:
                continue

            newest_version = version["tag"]
            break

        # Handle multiple instances of the same module
        for instance in installed[module["source"]]:
            try:
                cur = semver.VersionInfo.parse(instance["version"])
            except:
                # skip installed instanced with dev version
                continue

            # Version from remote repositories are already sorted
            # First match is the newest release
            if v > cur:
                # Create a copy to make to not edit original object
                update = instance.copy()
                update["update"] = version["tag"]
                updates.append(update)

    return updates


