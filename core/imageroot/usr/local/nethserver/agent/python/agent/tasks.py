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

import aiohttp
import asyncio
import json
import os
import urllib

async def _alogin(client, **npargs):
    """Read the Redis user credentials from the environment and retrieve a new authorization token.
    """
    tls_verify = npargs['tls_verify']
    endpoint = npargs['endpoint']
    redis_username = os.getenv('REDIS_USER', os.getenv('AGENT_ID', 'default'))
    redis_password = os.environ['REDIS_PASSWORD'] # Fatal if missing!
    async with client.post(
        f'{endpoint}/api/login', 
        json={
            'username':redis_username,
            'password':redis_password
        },
        ssl=tls_verify,
    ) as resp:
        jresp = await resp.json()

    if 'AGENT_STATE_DIR' in os.environ:
        cachefile = os.environ['AGENT_STATE_DIR'] + '/apitoken.cache'
        oldmask = os.umask(0o77)
        with open(cachefile, 'w') as fo:
            fo.write(jresp['token'])
        os.umask(oldmask)

    return {"Authorization": "Bearer " + jresp['token']}

async def _aread_status(client, agent_id, task_id, **npargs):
    """Read the status of the given task_id from agent_id
    """
    theaders = npargs['theaders']
    tls_verify = npargs['tls_verify']
    endpoint = npargs['endpoint']

    if theaders is None:
        theaders = _get_theaders_cache()

    async with client.get(
        f'{endpoint}/api/{agent_id}/task/{task_id}/status', 
        ssl=tls_verify,
        headers=theaders,
    ) as resp:
        jresp = await resp.json()

    return {
        'id': task_id,
        'output': jresp['data']['output'],
        'error': jresp['data']['error'],
        'exit_code': jresp['data']['exit_code'],
    }

def _get_theaders_cache():
    """Build authorization headers from the file cache contents.
    """
    if 'AGENT_STATE_DIR' in os.environ:
        cachefile = os.environ['AGENT_STATE_DIR'] + '/apitoken.cache' 
    else:
        cachefile = None

    theaders = {}
    try:
        with open(cachefile) as f:
            token = f.read().strip()
            theaders["Authorization"] = "Bearer " + token
    except:
        pass

    return theaders

async def _apost_task(client, agent_id, action, data, parent, **npargs):
    """Create a new task for the given agent_id and return its task_id.
    """
    theaders = npargs['theaders']
    tls_verify = npargs['tls_verify']
    endpoint = npargs['endpoint']
    if theaders is None:
        theaders = _get_theaders_cache()

    async with client.post(
        f'{endpoint}/api/{agent_id}/tasks',
        json={
            'action': action,
            'data': data,
            'parent': parent,
        },
        ssl=tls_verify,
        headers=theaders,
    ) as resp:

        jresp = await resp.json()

    return jresp['data']['id']

def _get_token(theaders):
    if not 'Authorization' in theaders:
        return 'eyJhbGciOiJIUzI1NiJ9.e30.qvi3rANM-kA4AAqkds0uf7bMgM7bpsYsCJWFoXKX9nQ' # an invalid token
    return theaders['Authorization'][len('Bearer '):]

async def _amonitor_task(client, agent_id, action, data, parent, **npargs):
    """Run the coro_post coroutine, expecting the task_id to watch from its return value.
    The coroutine is started after the web socket is established, thus any message related to task_id are
    caught.
    """
    task_id = None
    theaders = npargs['theaders']
    tls_verify = npargs['tls_verify']
    endpoint = npargs['endpoint']
    async with client.ws_connect(
        f'{endpoint}/ws?jwt=' + urllib.parse.quote_plus(_get_token(theaders)),
        ssl=tls_verify,
        heartbeat=4.0,
    ) as ws:
        task_id = await _apost_task(client, agent_id, action, data, parent, **npargs)
        async for msg in ws:
            if msg.type != aiohttp.WSMsgType.TEXT:
                continue

            try:
                jdata=json.loads(msg.data)
            except:
                continue # Discard unknown message

            if jdata['payload']['status'] != 'completed':
                continue # Ignore any status but "completed" XXX: catch "progress" too!

            if task_id and jdata['name'] == f'progress/{agent_id}/task/{task_id}':
                break

    return task_id

def run(agent_id, action, data={}, parent=None, progress_range=None, tls_verify=False, endpoint='http://cluster-leader/cluster-admin'):
    """Run a new task and wait until it completes. An object with exit_code, error, output is returned.
    """

    if parent is None:
        parent = os.getenv("AGENT_TASK_ID", "")

    # There are two main tasks that run concurrently. We must ensure the web socket captures any status update
    # of the Task even before it is created. This ensures we can catch the "completed" message and retrieve the
    # Task status from Redis after it becomes available.
    async def _asyncf():
        nonlocal parent
        theaders = _get_theaders_cache()
        pconn = {
            'tls_verify': tls_verify,
            'endpoint': endpoint,
        }
        async with aiohttp.ClientSession(raise_for_status=True) as client:
            try:
                task_id = await _amonitor_task(client, agent_id, action, data, parent, theaders=theaders, **pconn)
            except aiohttp.ClientResponseError as ex:
                if not ex.status in [200, 401]: # WSServerHandshakeError has status 200
                    raise ex
                theaders = await _alogin(client, **pconn)
                task_id = await _amonitor_task(client, agent_id, action, data, parent, theaders=theaders, **pconn)

            try:
                return await _aread_status(client, agent_id, task_id, theaders=theaders, **pconn)
            except aiohttp.ClientResponseError as ex:
                if ex.status != 401:
                    raise ex
                theaders = await _alogin(client, **pconn)
                return await _aread_status(client, agent_id, task_id, theaders=theaders, **pconn)

    return asyncio.run(_asyncf())
