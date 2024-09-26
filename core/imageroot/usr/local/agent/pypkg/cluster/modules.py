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

_repo_testing_cache = {}
def _repo_has_testing_flag(rdb, repo_name):
    """Retrieve from Redis the testing flag of repo_name one time, then read it from a cache."""
    global _repo_testing_cache
    if repo_name not in _repo_testing_cache:
        _repo_testing_cache[repo_name] = rdb.hget(f'cluster/repository/{repo_name}', 'testing') == "1"
    return _repo_testing_cache[repo_name]

_repo_view = None
def select_repo_view(view_name):
    """Change the module global configuration to send HTTP requests for a
    specific view to the metadata HTTPS service. Returned contents depend
    on the server side policy. By default, the sent value is "latest": use
    this function to change it, but do not call it more than one time in
    the same Python script and call it before any other function of this
    Python package."""
    global _repo_view
    if _repo_view is None:
        _repo_view = view_name
    else:
        raise Exception(f"View {_repo_view} is already selected. The view cannot be selected more than one time in the same script.")

def get_repo_view():
    """Return the selected the repository view. See
    selected_repo_view()."""
    global _repo_view
    if _repo_view is None:
        return 'latest'
    return _repo_view

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

def _calc_certification_level(package, has_subscription=False):
    if package['repository_authority'] in ["distfeed.nethserver.org", "subscription.nethserver.com"]:
        certification_level = 3
    elif package['repository_authority'] == "forge.nethserver.org":
        certification_level = 2
    else:
        certification_level = 1
    # If we trust the repo metadata, elevate up to level 5
    if certification_level == 3:
        if package["source"].startswith("ghcr.io/nethserver/") or package["source"].startswith("ghcr.io/nethesis/"):
            certification_level = 5 if has_subscription else 4
    return certification_level

def _parse_repository_metadata(repository_name, repository_url, repository_updated, metadata):
    try:
        repodata = json.loads(metadata)
    except Exception as ex:
        print(agent.SD_WARNING + f"Unable to parse metadata of repository {repository_name} at {repository_url}:", ex, file=sys.stderr)
        return []
    modules = []
    for package in repodata:
        package["repository"] = repository_name
        package["repository_updated"] = repository_updated
        try:
            package['repository_authority'] = urllib.parse.urlparse(repository_url).hostname
        except Exception as ex:
            package['repository_authority'] = ""
            print(agent.SD_WARNING + f"Unable to parse repository {repository_name} URL: {repository_url}", ex, file=sys.stderr)
        # Set absolute path for logo
        if package["logo"]:
            package["logo"] = _urljoin(repository_url, package["id"], package["logo"])
        else:
            package["logo"] = "" # convert None to empty string
        # Set absolute path for screenshots
        screenshots = []
        for s in package["screenshots"]:
           screenshots.append(_urljoin(repository_url, package["id"], s))
        package["screenshots"] = screenshots
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

def _list_repository_modules(rdb, repository_name, repository_url):
    cache_key = f'cluster/repository_cache/{repository_name}'
    hcache = rdb.hgetall(cache_key)
    repo_view = get_repo_view()
    if not hcache or hcache.get('repo_view') != repo_view:
        url = _urljoin(repository_url, "repodata.json")
        hsubscription = rdb.hgetall("cluster/subscription") or None
        try:
            with _get_http_session() as osession:
                if hsubscription and url.startswith("https://subscription.nethserver.com/"):
                    # Send system_id for HTTP Basic authentication
                    osession.auth = (hsubscription["system_id"], hashlib.sha256(hsubscription["auth_token"].encode()).hexdigest())
                resp = osession.get(url, params={"view": repo_view})
                repodata_raw = resp.text
                updated = resp.headers.get('Last-Modified', "")
        except Exception as ex:
            print(f"Fetching {url}:", ex, file=sys.stderr)
            # If repository is not accessible or invalid, just return an empty array
            return []
        hcache = {"data": repodata_raw, "updated": updated, "repo_view": repo_view}
    modules = _parse_repository_metadata(repository_name, repository_url, hcache['updated'], hcache['data'])
    # Save inside the cache if data is valid
    if modules:
        # Save also repodata file date
        rdb.hset(cache_key, mapping=hcache)
        # Set cache expiration to 3600 seconds
        rdb.expire(cache_key, 3600)
    return modules

class LatestModuleLookupError(Exception):
    pass

def get_latest_module(module, rdb):
    """Find most recent version of the given module
    """
    for m in list_available(rdb, skip_core_modules=False):
        if m["id"] == module:
            # We assume at index 0 we find the latest tag:
            return m["source"] + ':' + m['versions'][0]['tag']

    raise LatestModuleLookupError(module)

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
    ometadata.setdefault("repository_authority", "__local__")
    ometadata["certification_level"] = 0
    ometadata["rootfull"] = False
    try:
        ometadata['logo'] = glob(f'{path_prefix}apps/{module_id}/img/*logo*png')[0].removeprefix(path_prefix)
    except Exception as ex:
        ometadata['logo'] = ""
        print(agent.SD_INFO + "_fetch_metadata_json/glob:", ex, file=sys.stderr)
    return ometadata

def _tag_is_stable(tag):
    """Return true if the tag string represents a Semver stable release"""
    try:
        return not bool(semver.Version.parse(tag).prerelease)
    except:
        pass
    return False

def _min_core_ok(modver, leader_core_version):
    try:
        min_core = semver.Version.parse(modver['labels']['org.nethserver.min-core'])
    except:
        min_core = semver.Version(0)
    return leader_core_version >= min_core

def list_available(rdb, skip_core_modules=False):
    """Iterate over enabled repositories and return available modules respecting the repository priority."""
    hsubscription = rdb.hgetall("cluster/subscription") or None
    leader_core_version = _get_leader_core_version()
    modules = []
    for omod in _get_available_modules(rdb).values():
        if not _repo_has_testing_flag(rdb, omod["repository"]):
            # Ignore testing releases for new installations:
            omod["versions"] = list(filter(lambda v: _tag_is_stable(v["tag"]), omod["versions"]))
        omod["versions"] = list(filter(lambda v: _min_core_ok(v, leader_core_version), omod["versions"]))
        if not omod["versions"]:
            continue # Ignore modules with no versions
        omod["certification_level"] = _calc_certification_level(omod, bool(hsubscription))
        try:
            if skip_core_modules and 'core_module' in omod["versions"][0]['labels']['org.nethserver.flags']:
                continue # core modules are ignored
        except:
            pass
        try:
            package_is_rootfull = omod["versions"][0]["labels"]["org.nethserver.rootfull"] == "1"
        except:
            package_is_rootfull = False
        # Ignore untrusted rootfull application, if a subscription is active
        if hsubscription and package_is_rootfull and omod["certification_level"] < 3:
            print(agent.SD_WARNING + f"Ignoring image of rootfull application {omod['source']}: certification_level {omod['certification_level']} is too low", file=sys.stderr)
            continue # skip package
        else:
            omod['rootfull'] = package_is_rootfull
        modules.append(omod)
    return modules

def _get_available_modules(rdb):
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
        for rmod in _list_repository_modules(rdb, nrepo, repo["url"]):
            if rmod["source"] in modules:
                continue # skip duplicated images from lower priority modules
            modules[rmod["source"]] = rmod
            rmod['versions'].sort(key=lambda v: _parse_version_object(v["tag"]), reverse=True)
    # Integrate the available set with instances that do not belong to any
    # repository. They can be found in the "installed" dict:
    for module_source, module_instances in list_installed(rdb).items():
        if module_source in modules:
            continue
        _, image_name = module_source.rsplit("/", 1)
        vmetadata = _fetch_metadata_json(module_instances[0]['id'], image_name)
        vmetadata["versions"] = list(_synthesize_module_version(oinst) for oinst in module_instances)
        vmetadata["versions"].sort(key=lambda v: _parse_version_object(v["tag"]), reverse=True)
        vmetadata["source"] = module_source
        vmetadata["id"] = image_name
        modules[module_source] = vmetadata
    return modules

def list_installed(rdb, skip_core_modules = False):
    installed = {}
    logos = _get_downloaded_logos()
    hmodules = rdb.hgetall("cluster/module_node") or {}
    hnode_names = {}
    for node_id in set(hmodules.values()):
        hnode_names[node_id] = rdb.get(f'node/{node_id}/ui_name') or ""
    for module_id in hmodules.keys():
        try:
            mflags = list(rdb.smembers(f'module/{module_id}/flags'))
            if skip_core_modules and 'core_module' in mflags:
                continue # ignore core modules as requested
            menv = rdb.hgetall(f'module/{module_id}/environment') or {}
            msource, mtag = menv['IMAGE_URL'].rsplit(":", 1)
            mnode_id = hmodules[module_id]
            hinstance = {
                'id': module_id,
                'source': msource,
                'version': mtag,
                'module': msource.rsplit("/", 1)[1],
                'ui_name': rdb.get(f'module/{module_id}/ui_name') or "",
                'node': mnode_id,
                'node_ui_name': hnode_names[mnode_id],
                'logo': logos.get(module_id, ""),
                'digest': menv['IMAGE_DIGEST'],
                'flags': mflags,
            }
            installed.setdefault(msource, [])
            installed[msource].append(hinstance)
        except Exception as ex:
            print(agent.SD_ERR+f"Cannot fetch {module_id} attributes:", ex, file=sys.stderr)
            continue
    for instances in installed.values():
        instances.sort(key=lambda v: _parse_version_object(v["version"]), reverse=True)
    return installed

def _get_core_tag():
    core_env = agent.read_envfile('/etc/nethserver/core.env')
    _, tag = core_env['CORE_IMAGE'].rsplit(":", 1)
    return tag

def _get_leader_core_version():
    try:
        leader_core_version = semver.parse_version_info(_get_core_tag())
    except:
        leader_core_version = semver.Version(999)
    return leader_core_version

def list_updates(rdb, skip_core_modules=False, with_testing_update=False):
    updates = []
    installed_modules = list_installed(rdb, skip_core_modules)
    available_modules = _get_available_modules(rdb)
    leader_core_version = _get_leader_core_version()

    node_core_versions = _get_node_core_versions(rdb)
    flat_instance_list = list(mi for module_instances in installed_modules.values() for mi in module_instances)
    for instance in flat_instance_list:
        try:
            current_core = semver.parse_version_info(node_core_versions.get(instance["node"]))
        except:
            current_core = leader_core_version
        if not instance['source'] in available_modules:
            continue # skip instance if is not available from any repository
        try:
            current_version = semver.parse_version_info(instance['version'])
        except:
            continue # skip development version: instance must be updated manually
        # Assuming the versions array is sorted in decreasing oreder, look
        # up update candidates for both stable and testing updates:
        update_candidate = None
        testing_update_candidate = None
        available_module = available_modules[instance['source']]
        repository_name = available_module['repository']
        for aver in available_module['versions']:
            try:
                available_version = semver.parse_version_info(aver['tag'])
            except:
                continue # skip non-semver available tag
            if available_version <= current_version:
                continue # ignore tags that do not update the current one
            try:
                minimum_version = semver.parse_version_info(aver['labels']['org.nethserver.min-from'])
            except:
                # Arbitrary low version to satisfy any tag:
                minimum_version = (0,0,0)
            if current_version < minimum_version:
                print(agent.SD_NOTICE + f"Ignoring update of {instance['id']} with {instance['source']}:{aver['tag']}: org.nethserver.min-from", minimum_version, file=sys.stderr)
                continue # Skip versions incompatible with instance version.
            try:
                minimum_core = semver.parse_version_info(aver['labels']['org.nethserver.min-core'])
            except:
                minimum_core = (0,0,0)
            if current_core < minimum_core:
                print(agent.SD_NOTICE + f"Ignoring update of {instance['id']} with {instance['source']}:{aver['tag']}: org.nethserver.min-core:", minimum_core, file=sys.stderr)
                continue # Skip versions incompatible with core.
            if update_candidate is None and (
                _repo_has_testing_flag(rdb, repository_name)
                or not available_version.prerelease
            ):
                update_candidate = available_version
            if testing_update_candidate is None and with_testing_update and available_version.prerelease:
                testing_update_candidate = available_version

        # If a stable or testing candidate has been found, add this
        # instance to the updates list.
        if update_candidate:
            instance['update'] = str(update_candidate)
        if testing_update_candidate:
            instance['testing_update'] = str(testing_update_candidate)
        if 'update' in instance or 'testing_update' in instance:
            updates.append(instance)

    return updates

def _get_node_core_versions(rdb):
    hversions = {}
    for node_id in set(rdb.hvals("cluster/module_node")):
        image_url = rdb.hget(f'node/{node_id}/environment', 'IMAGE_URL')
        if not image_url:
            print(agent.SD_WARNING+f"Cannot fetch IMAGE_URL of node {node_id}", file=sys.stderr)
            if node_id == os.getenv("NODE_ID"):
                image_url = "ghcr.io/nethserver/core:" + _get_core_tag()
            else:
                image_url = "ghcr.io/nethserver/core:0.0.0"
        _, vtag = image_url.rsplit(":", 1)
        hversions[node_id] = vtag
    return hversions

def get_node_core_versions(rdb):
    return _get_node_core_versions(rdb)

def list_core_modules(rdb):
    """List core modules and if they can be updated."""
    updates = list_updates(rdb, skip_core_modules=False)
    def _get_module_update(module_id):
        for oupdate in updates:
            if oupdate['id'] == module_id:
                return oupdate['update']
        return ""
    core_modules = {}
    try:
        _, latest_core = get_latest_module('core', rdb).rsplit(":", 1)
    except LatestModuleLookupError:
        latest_core = '0.0.0'
    def _calc_core_update(current, latest):
        try:
            vcur = semver.parse_version_info(current)
        except:
            vcur = semver.Version(999)
        try:
            vlatest = semver.parse_version_info(latest)
        except:
            vlatest = semver.Version(0)
        if vlatest > vcur:
            return latest
        else:
            return ""
    for node_id, ntag in _get_node_core_versions(rdb).items():
        try:
            core_instance = {
                'id': 'core' + node_id,
                'version': ntag,
                'update': _calc_core_update(ntag, latest_core),
                'node_id': node_id,
                'node_ui_name': rdb.get(f'node/{node_id}/ui_name') or "",
            }
        except Exception as ex:
            print(agent.SD_ERR+f"Cannot fetch node {node_id} attributes:", ex, file=sys.stderr)
            continue
        core_modules.setdefault('core', {"name": 'core', "instances": []})
        core_modules['core']['instances'].append(core_instance)
    for module_source, instances in list_installed(rdb, skip_core_modules=False).items():
        _, image_name = module_source.rsplit("/", 1)
        try:
            has_core_module = 'core_module' in instances[0]['flags']
        except:
            has_core_module = False
        if not has_core_module:
            continue
        core_modules.setdefault(image_name, {"name": image_name, "instances": []})
        for instance in instances:
            core_modules[image_name]['instances'].append({
                "id": instance["id"],
                "version": instance["version"],
                "update": _get_module_update(instance['id']),
                "node_id": instance["node"],
                "node_ui_name": instance["node_ui_name"],
            })
    return list(core_modules.values())
