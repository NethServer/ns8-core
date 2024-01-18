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

import redis
import ldap3
import ssl
import sys

def probe_ldap_basedn(config, ldapconn=None):
    if ldapconn is None:
        ldapconn = _open_ldap_connection(config, get_info=ldap3.DSA)
    try:
        return ldapconn.server.info.naming_contexts[0]
    except:
        return 'dc=' + config["domain"].lower().replace(".", ",dc=")

def probe_ldap_schema(config):
    ldapconn = _open_ldap_connection(config, get_info=ldap3.DSA)
    try:
        if ldapconn.server.info.other['isGlobalCatalogReady'] == ['TRUE']:
            return "ad"
    except:
        pass

    return "rfc2307"

def _open_ldap_connection(config, get_info=None):
    tls_verify = ssl.CERT_REQUIRED if config['tls_verify'] else ssl.CERT_NONE
    ldapsrv = ldap3.Server(
        config['host'],
        port=config['port'],
        use_ssl=config['tls'],
        tls=ldap3.Tls(validate=tls_verify) if config['tls'] else None,
        connect_timeout=5,
        get_info=get_info,
    )
    ldapconn = ldap3.Connection(ldapsrv,
        user=config['bind_dn'] or None,
        password=config['bind_password'] or None,
        client_strategy=ldap3.SAFE_SYNC,
        receive_timeout=5,
        auto_bind=True,
        raise_exceptions=True,
    )
    return ldapconn

def validate_ldap(config):
    logex = None
    errors = []
    ldapconn = None
    try:
        ldapconn = _open_ldap_connection(config, get_info=ldap3.DSA)
    except ldap3.core.exceptions.LDAPSocketOpenError as ex:
        logex = ex
        if('socket ssl wrapping' in str(ex)):
            errors.append({'field':'tls_verify', 'parameter':'tls_verify','value': config['tls_verify'], 'error': 'invalid_tls_certificate'})
        elif(str(ex) == 'invalid server address'):
            errors.append({'field':'host', 'parameter':'host','value': config['host'], 'error': 'invalid_server_address'})
        else:
            errors.append({'field':'port', 'parameter':'port','value': config['port'], 'error': 'port_connection_error'})
    except ldap3.core.exceptions.LDAPSocketReceiveError as ex:
        logex = ex
        errors.append({'field':'tls', 'parameter':'tls','value': config['tls'], 'error': 'wrong_service_port_or_tls_setting'})
        errors.append({'field':'port', 'parameter':'port','value': config['port'], 'error': 'wrong_service_port_or_tls_setting'})
    except ldap3.core.exceptions.LDAPInvalidCredentialsResult as ex:
        logex = ex
        errors.append({'field':'bind_dn', 'parameter':'bind_dn','value': config['bind_dn'], 'error': 'invalid_credentials'})
        errors.append({'field':'bind_password', 'parameter':'bind_password','value': None, 'error': 'invalid_credentials'})
    except ldap3.core.exceptions.LDAPStrongerAuthRequiredResult as ex:
        logex = ex
        errors.append({'field':'tls', 'parameter':'tls','value': config['tls'], 'error': 'tls_required_for_authentication'})

    except Exception as ex:
        logex = ex
        errors.append({'field':'host', 'parameter':'host','value': config['host'], 'error': 'generic_host_connection_error'})

    if ldapconn:
        if config['base_dn'] == "":
            base_dn = probe_ldap_basedn(config, ldapconn)
        else:
            base_dn = config['base_dn']

        try:
            ldapconn.search(base_dn, '(objectClass=*)', search_scope=ldap3.BASE)
        except ldap3.core.exceptions.LDAPNoSuchObjectResult as ex:
            logex = ex
            errors.append({'field':'base_dn', 'parameter':'base_dn','value': base_dn, 'error': 'base_dn_not_found'})
        except ldap3.core.exceptions.LDAPInvalidDnError as ex:
            logex = ex
            errors.append({'field':'base_dn', 'parameter':'base_dn','value': base_dn, 'error': 'dn_syntax_error'})
        except Exception as ex:
            logex = ex
            errors.append({'field':'base_dn', 'parameter':'base_dn','value': base_dn, 'error': 'generic_base_dn_check_error'})

    return (errors, logex)

def get_internal_domains(rdb):
    """Read from Redis the internal user domains configuration
    """
    domains={}
    configured_providers = set()

    for kldap in rdb.scan_iter("module/*/srv/tcp/ldap"):
        module_id = kldap.split('/', 2)[1]
        configured_providers.add(module_id)
        conf = rdb.hgetall(kldap)

        if not conf['domain'] in domains:
            domains[conf['domain']] = {
                "name": conf['domain'],
                "location": "internal",
                "protocol": "ldap",
                "schema": conf['schema'],
                "tls": conf['tls'] == 'on',
                "tls_verify": conf['tls_verify'] == 'on',
                "base_dn": conf['base_dn'],
                "bind_dn": conf['bind_dn'],
                "bind_password": conf['bind_password'],
                "providers": []
            }

        # Check if the provider can be used also as a SMB file server by
        # LAN clients. This is possible only with Samba providers but the
        # attribute is always set despite of that.
        has_file_server_flag = bool(rdb.sismember(f'module/{module_id}/flags', 'file_server'))

        domains[conf['domain']]['providers'].append({
            "id": module_id,
            "ui_name": rdb.get(f'module/{module_id}/ui_name') or "",
            "node": int(conf['node']),
            "file_server": has_file_server_flag,
            "host": conf['host'],
            "port": int(conf['port']),
        })

    for kud in rdb.scan_iter("module/*/user_domain"):
        module_id = kud.split('/', 2)[1]
        if module_id in configured_providers:
            continue

        domain_name = rdb.get(kud)
        if not domain_name in domains:
            continue # skip new domain module instance

        node_id = rdb.hget(f'module/{module_id}/environment', 'NODE_ID')
        domains[domain_name]['providers'].append({
            "id": module_id,
            "ui_name": rdb.get(f'module/{module_id}/ui_name') or "",
            "node": int(node_id),
            "host": None,
            "port": None,
            "file_server": False, # An unconfigured internal provider cannot be a file server
        })

    return domains

def get_external_domains(rdb):
    """Read from Redis the external user domains configuration
    """
    domains = {}
    for kldap in rdb.scan_iter("cluster/user_domain/ldap/*/conf"):
        domain_id = kldap.rsplit('/', 2)[1]
        conf = rdb.hgetall(kldap)

        domains[domain_id] = {
            "name": domain_id,
            "location": "external",
            "protocol": "ldap",
            "schema": conf['schema'],
            "tls": conf['tls'] == 'on',
            "tls_verify": conf['tls_verify'] == 'on',
            "base_dn": conf['base_dn'],
            "bind_dn": conf['bind_dn'],
            "bind_password": conf['bind_password'],
            "providers": []
        }

        providers = rdb.lrange(f"cluster/user_domain/ldap/{domain_id}/providers", 0, 8)
        for epaddr in providers:
            host, port = epaddr.split(':', 1)
            domains[domain_id]['providers'].append({
                "host": host,
                "port": int(port),
                "id": host,
                "file_server": False, # An external provider cannot be a file server
                "node": None,
                "ui_name": rdb.hget(f"cluster/user_domain/ldap/{domain_id}/ui_names", epaddr) or "",
            })

        if not domains[domain_id]['providers']:
            del domains[domain_id]

    return domains

def list_domains(rdb):
    # Internal domains have higher precedence and override external ones:
    domains = get_external_domains(rdb)
    domains.update(get_internal_domains(rdb))
    return domains

def get_domain_modules(rdb, domain):
    """Return a list of module identifiers that are bound to the given domain"""

    module_list = []
    rawrel = rdb.hgetall("cluster/module_domains") or {}
    for (key, val) in rawrel.items():
        pdomains = val.split()
        if domain in pdomains:
            module_list.append(key)

    return module_list
