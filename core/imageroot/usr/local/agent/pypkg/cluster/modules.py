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
import urllib
import requests, urllib3.util
import hashlib
from glob import glob
import sys
import subprocess
import datetime

def _urljoin(base_path, *args):
    '''replace urllib.parse.joinurl because it doesn't handle multiple parameters
    '''
    parts = [base_path]
    for arg in args:
        # also make sure to escape special chars using quote
        parts.append(urllib.parse.quote(arg))
    return "/".join(part.strip("/") for part in parts)

def _get_downloaded_logos():
    logos = {}
    # Search for log inside all caches
    for app in glob("/var/lib/nethserver/cluster/ui/apps/*"):
        src_logos = glob(f'{app}/img/*logo*png')
        if len(src_logos) > 0:
            logo = src_logos[0].removeprefix("/var/lib/nethserver/cluster/ui/")
            logos[os.path.basename(app)] = logo
    return logos

def _parse_repository_metadata(repository_name, repository_url, repository_updated, repodata, skip_core_modules = False, skip_testing_versions = False):
    modules = []

    try:
        repodata = json.loads(repodata)
    except:
        return modules

    def ignore_testing(version):
        if skip_testing_versions and version["testing"] is True:
            return False
        else:
            return True

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

        # Filter
        package["versions"] = list(filter(ignore_testing, package["versions"]))

        if len(package["versions"]) > 0:
            modules.append(package)

    return modules

def _get_http_session():
    osession = requests.Session()
    osession.headers= {
        'User-Agent': 'ns8-downloader', # DO Spaces blocks Python UA
    }
    osession.timeout = 15 # Timeout for HTTP connections
    oretries = urllib3.util.Retry(
        total=3,
        backoff_factor=0.1,
        allowed_methods={'POST', 'GET'},
    )
    osession.mount('https://', requests.adapters.HTTPAdapter(max_retries=oretries))
    return osession

def _list_repository_modules(rdb, repository_name, repository_url, skip_core_modules = False, skip_testing_versions=False):
    key = f'cluster/repository_cache/{repository_name}'
    cache = rdb.hgetall(key)
    if cache:
        return _parse_repository_metadata(repository_name, repository_url, cache["updated"], cache["data"], skip_core_modules, skip_testing_versions)

    url = _urljoin(repository_url, "repodata.json")
    try:
        hsubscription = rdb.hgetall("cluster/subscription")
        with _get_http_session() as osession:
            if hsubscription and url.startswith("https://subscription.nethserver.com/"):
                # Send system_id for HTTP Basic authentication
                osession.auth = (hsubscription["system_id"], hashlib.sha256(hsubscription["auth_token"].encode()).hexdigest())
            resp = osession.get(url)
            repodata = resp.text
            updated = resp.headers.get('Last-Modified', "")
    except Exception as ex:
        print(f"Fetching {url}:", ex, file=sys.stderr)
        # If repository is not accessible or invalid, just return an empty array
        return []

    modules = _parse_repository_metadata(repository_name, repository_url, updated, repodata, skip_core_modules, skip_testing_versions)
    # Save inside the cache if data is valid
    if modules:
        # Save also repodata file date
        rdb.hset(key, mapping={"data": repodata, "updated": updated})
        # Set cache expiration to 3600 seconds
        rdb.expire(key, 3600)

    return modules

class LatestModuleLookupError(Exception):
    pass

def get_latest_module(module, rdb):
    """Find most recent version of the given module
    """
    version = ""
    source = ""
    available = list_available(rdb)
    for m in available:
        if m["id"] == module:
            source = m["source"]
            repo_testing = rdb.hget(f'cluster/repository/{m["repository"]}', 'testing') == "1"
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
        raise LatestModuleLookupError(module)

    return f'{source}:{version}'

def _parse_version_object(v):
    try:
        vinfo = semver.Version.parse(v)
    except:
        vinfo = semver.Version(0)
    return vinfo

def _synthesize_module_version(minstance):
    """From an installed module instance, synthesize module version for
    repository metadata, or fall back to a safe default."""
    inspect_labels = {}
    image_url = minstance["source"] + ':' + minstance["version"]
    # Extract image labels
    try:
        with subprocess.Popen(['podman', 'image', 'inspect', image_url], stdout=subprocess.PIPE, stderr=sys.stderr) as proc:
            inspect = json.load(proc.stdout)
            inspect_labels = inspect[0]['Labels']
    except Exception as ex:
        print(agent.SD_WARNING + "_synthesize_module_version:", ex, file=sys.stderr)
    # Extract the prerelease flag:
    try:
        testing = bool(semver.Version.parse(minstance["version"]).prerelease)
    except:
        testing = True
    inspect_labels.setdefault("org.nethserver.images", image_url)
    return {
        "tag": minstance["version"],
        "testing": testing,
        "labels": inspect_labels,
    }

def _fetch_metadata_json(module_id, image_name):
    """Synthesize repository metadata from the local metadata.json
    file, or fall back to safe defaults."""
    path_prefix = '/var/lib/nethserver/cluster/ui/'
    repository_updated_timestamp = 'Thu Jun 13 01:00:00 AM UTC 2024'
    try:
        with open(f"{path_prefix}apps/{module_id}/metadata.json") as mfd:
            # Read modification timestamp of the .json file:
            repository_updated_timestamp = datetime.datetime.utcfromtimestamp(os.fstat(mfd.fileno()).st_mtime).ctime()
            ometadata = json.load(mfd)
    except Exception as ex:
        print(agent.SD_INFO + "_fetch_metadata_json/open:", ex, file=sys.stderr)
        ometadata = {
            "id": image_name,
            "name": image_name.capitalize(),
            "description": {"en": ""},
            "categories": [],
            "logo": "",
            "authors": [],
            "docs": {
                "documentation_url": "",
                "code_url": "",
                "bug_url": "",
            },
        }
    ometadata.setdefault("screenshots", [])
    ometadata.setdefault("repository", "__local__")
    ometadata.setdefault("repository_updated", repository_updated_timestamp)
    try:
        ometadata['logo'] = glob(f'{path_prefix}apps/{module_id}/img/*logo*png')[0].removeprefix(path_prefix)
    except Exception as ex:
        ometadata['logo'] = ""
        print(agent.SD_INFO + "_fetch_metadata_json/glob:", ex, file=sys.stderr)
    return ometadata

def list_available(rdb, skip_core_modules = False):
    """Iterate over enabled repositories and return available modules respecting the repository priority."""
    modules = {}
    repositories = []
    # List all modules from enabled repositories
    for krepo in rdb.scan_iter('cluster/repository/*'):
        repositories.append(krepo.removeprefix("cluster/repository/"))
    # Alphabetical order, where last item has higher priority:
    repositories.sort(reverse=True)
    for nrepo in repositories:
        repo = rdb.hgetall('cluster/repository/' + nrepo)
        # Skip non-enabled repositories
        if repo.get("status", "0") != "1":
            continue
        skip_testing_versions = repo.get("testing", "0") != "1"
        for rmod in _list_repository_modules(rdb, nrepo, repo["url"], skip_core_modules, skip_testing_versions):
            if rmod["source"] in modules:
                continue # skip duplicated images from lower priority modules
            modules[rmod["source"]] = rmod
            rmod['versions'].sort(key=lambda v: _parse_version_object(v["tag"]), reverse=True)
    # Integrate the available set with instances that do not belong to any
    # repository. They can be found in the "installed" dict:
    for module_source, module_instances in list_installed(rdb, skip_core_modules).items():
        if module_source in modules:
            continue
        _, image_name = module_source.rsplit("/", 1)
        vmetadata = _fetch_metadata_json(module_instances[0]['id'], image_name)
        vmetadata["versions"] = list(_synthesize_module_version(oinst) for oinst in module_instances)
        vmetadata["installed"] = module_instances
        vmetadata["source"] = module_source
        vmetadata["updates"] = []
        vmetadata["id"] = image_name
        modules[module_source] = vmetadata
    return list(modules.values())

def list_installed(rdb, skip_core_modules = False):
    installed = {}
    logos = _get_downloaded_logos()
    # Search for installed modules
    for m in rdb.scan_iter('module/*/environment'):
        vars = rdb.hgetall(m)
        module_ui_name = rdb.get(m.removesuffix('/environment') + '/ui_name') or ""
        url, sep, tag = vars['IMAGE_URL'].partition(":")
        image = url[url.rindex('/')+1:]
        logo = logos.get(vars["MODULE_ID"]) or ''
        flags = list(rdb.smembers(f'module/{vars["MODULE_ID"]}/flags')) or []
        if skip_core_modules and 'core_module' in flags:
            continue
        if url not in installed.keys():
            installed[url] = []
        # Retrieve node ui_name
        node_id = vars['NODE_ID']
        node_ui_name = rdb.get(f"node/{node_id}/ui_name") or ""
        installed[url].append({ 'id': vars["MODULE_ID"], 'ui_name': module_ui_name, 'node': node_id, 'node_ui_name': node_ui_name, 'digest': vars["IMAGE_DIGEST"], 'source': url, 'version': tag, 'logo': logo, 'module': image, 'flags': flags})

    for instances in installed.values():
        instances.sort(key=lambda v: _parse_version_object(v["version"]), reverse=True)

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
            try:
                # skip bogus version tag
                v = semver.Version.parse(version["tag"])
            except:
                continue
            # Skip testing versions if testing is disabled
            testing = rdb.hget(f'cluster/repository/{module["repository"]}', 'testing')
            if testing != "1" and not v.prerelease is None:
                continue
            newest_version = version["tag"]
            break

        # Handle multiple instances of the same module
        for instance in installed[module["source"]]:
            try:
                cur = semver.Version.parse(instance["version"])
            except:
                # skip installed instanced with dev version
                continue

            # Version are already sorted
            # First match is the newest release
            if v > cur:
                # Create a copy to not change original object
                update = instance.copy()
                update["update"] = version["tag"]
                updates.append(update)

    return updates

def list_core_modules(rdb):
    """List core modules and if they can be updated."""
    core_modules = {}
    available = list_available(rdb, skip_core_modules = False)

    def _calc_update(image_name, cur):
        # Lookup module information from repositories
        for module in available:
            if module["id"] == image_name:
                break
        else:
            return ""
        try:
            vupdate = module["versions"][0]['tag']
            vinfo = semver.VersionInfo.parse(vupdate)
            if vupdate > cur:
                return vupdate
        except:
            pass
        return ""

    for module_source, instances in list_installed_core(rdb).items():
        _, image_name = module_source.rsplit("/", 1)
        core_modules.setdefault(image_name, {"name": image_name, "instances": []})
        for instance in instances:
            core_modules[image_name]['instances'].append({
                "id": instance["id"],
                "version": instance["version"],
                "update": _calc_update(image_name, instance["version"]),
            })

    return list(core_modules.values())
