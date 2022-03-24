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

def _urljoin(base_path, *args):
    '''replace urllib.parse.joinurl because it doesn't handle multiple parameters
    '''
    parts = [base_path]
    for arg in args:
        # also make sure to escape special chars using quote
        parts.append(urllib.parse.quote(arg))
    return "/".join(part.strip("/") for part in parts)

def _get_cached_logo_urls(rdb):
    logos = {}
    # Search for log inside all caches
    for m in rdb.scan_iter('cluster/repository_cache/*'):
        cache = rdb.hgetall(m)
        if cache:
            (prefix,sep,repo) = m.rpartition("/")
            repository_url = rdb.hget(f'cluster/repository/{repo}', 'url')
            for m in  _parse_repository_metadata(repo, repository_url, cache["updated"], cache["data"]):
                logos[m["id"]] = m['logo']
    return logos

def _parse_repository_metadata(repository_name, repository_url, repository_updated, repodata, skip_core_modules = False):
    modules = []

    try:
        repodata = json.loads(repodata)
    except:
        return modules

    for package in repodata:
        # Skip core modules if flag is enabled
        if skip_core_modules and package['versions']:
            version = package['versions'][0]
            if 'org.nethserver.flags' in version['labels'] and 'core_module' in version['labels']['org.nethserver.flags']:
                continue

        package["repository"] = repository_name
        package["repository_updated"] = repository_updated

        # Set absolute path for logo
        if package["logo"]:
           package["logo"] = _urljoin(repository_url, package["id"], package["logo"])

        # Set absolute path for screenshots
        screenshots = []
        for s in package["screenshots"]:
           screenshots.append(_urljoin(repository_url, package["id"], s))
        package["screenshots"] = screenshots

        modules.append(package)

    return modules


def _list_repository_modules(rdb, repository_name, repository_url, skip_core_modules = False):
    key = f'cluster/repository_cache/{repository_name}'
    cache = rdb.hgetall(key)
    if cache:
        return _parse_repository_metadata(repository_name, repository_url, cache["updated"], cache["data"], skip_core_modules)

    try:
        url = urllib.parse.urljoin(repository_url, "repodata.json")
        with urllib.request.urlopen(url) as req:
            repodata = req.read().decode()
    except:
        # If repository is not accessible or invalid, just return an empty array
        return []

    updated = req.headers.get('Last-Modified', "")
    modules = _parse_repository_metadata(repository_name, repository_url, updated, repodata, skip_core_modules)
    # Save inside the cache if data is valid
    if modules:
        # Save also repodata file date
        rdb.hset(key, mapping={"data": repodata, "updated": updated})
        # Set cache expiration to 3600 seconds
        rdb.expire(key, 3600)

    return modules

def get_latest_module(module, rdb):
    """Find most recent version of the given module
    """
    version = ""
    source = ""
    available = list_available(rdb)
    for m in available:
        if m["id"] == module:
            source = m["source"]
            repo_testing = int(rdb.hget(f'cluster/repository/{m["repository"]}', 'testing')) == 1
            for v in m["versions"]:
                if repo_testing:
                    version = v["tag"]
                elif not v["testing"]:
                    version = v["tag"]

                if version:
                    break

        if version:
            break


    # Fail if package has not been found inside the repository metadata
    if not source or not version:
        raise Exception(f'Package not found: {module}')

    return f'{source}:{version}'


def list_available(rdb, skip_core_modules = False):
    modules = []
    # List all modules from enabled repositories
    for m in rdb.scan_iter('cluster/repository/*'):
        repo = rdb.hgetall(m)
        # Skip disabled repositories
        if int(repo["status"]) == 0:
            continue
        
        modules.extend(_list_repository_modules(rdb, os.path.basename(m), repo["url"], skip_core_modules))

    return modules

def list_installed(rdb, skip_core_modules = False):
    installed = {}
    logos = _get_cached_logo_urls(rdb)
    # Search for installed modules
    for m in rdb.scan_iter('module/*/environment'):
        vars = rdb.hgetall(m)
        module_ui_name = rdb.get(m.removesuffix('/environment') + '/ui_name') or ""
        url, sep, tag = vars['IMAGE_URL'].partition(":")
        image = url[url.rindex('/')+1:]
        logo = logos.get(image) or ''
        flags = list(rdb.smembers(f'module/{vars["MODULE_ID"]}/flags')) or []
        if skip_core_modules and 'core_module' in flags:
            continue
        if url not in installed.keys():
            installed[url] = []
        installed[url].append({ 'id': vars["MODULE_ID"], 'ui_name': module_ui_name, 'node': vars['NODE_ID'], 'digest': vars["IMAGE_DIGEST"], 'source': url, 'version': tag, 'logo': logo, 'module': image, 'flags': flags})

    return installed

def list_installed_core(rdb):
    installed = {'ghcr.io/nethserver/core': []}
    core_env = agent.read_envfile('/etc/nethserver/core.env')
    (url, tag) = core_env['CORE_IMAGE'].split(":")
    installed['ghcr.io/nethserver/core'].append({ 'id': 'core', 'version': tag, 'module': 'core'})
    # Search for installed modules
    for m in rdb.scan_iter('module/*/environment'):
        vars = rdb.hgetall(m)
        if 'core_module' in rdb.smembers(f'module/{vars["MODULE_ID"]}/flags'):
            url, sep, tag = vars['IMAGE_URL'].partition(":")
            image = url[url.rindex('/')+1:]
        
            if url not in installed.keys():
                installed[url] = []
            installed[url].append({'id': vars["MODULE_ID"], 'version': tag, 'module': image})

    return installed


def list_updates(rdb, skip_core_modules = False):
    updates = []
    installed = list_installed(rdb, skip_core_modules)
    available = list_available(rdb, skip_core_modules)

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
                # Create a copy to not change original object
                update = instance.copy()
                update["update"] = version["tag"]
                updates.append(update)

    return updates

def list_core_modules():
    rdb = agent.redis_connect(privileged=True)
    modules = []
    installed = list_installed_core(rdb)
    available = list_available(rdb, skip_core_modules = False)

    for module in available:
        if not 'core' in module['categories']:
           continue

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
        instances = []
        for instance in installed[module["source"]]:
            del(instance['module'])
            instance["update"] = ""
            instances.append(instance)
            try:
                cur = semver.VersionInfo.parse(instance["version"])
            except:
                # skip installed instanced with dev version
                continue

            # Version from remote repositories are already sorted
            # First match is the newest release
            if v > cur:
                instance["update"] = version["tag"]

        modules.append({"name": module['name'], "instances": instances})

    return modules
