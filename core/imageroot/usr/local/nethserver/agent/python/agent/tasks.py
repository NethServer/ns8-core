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
import aioredis
import asyncio
import json
import os
import urllib
import sys
import uuid

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

def _get_theaders_cache(token=None):
    """Build authorization headers from the file cache contents.
    """
    if token:
        return {
            "Authorization": "Bearer " + token
        }

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
        return ''
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

            if not jdata['payload']['status'] in ['completed', 'aborted', 'validation-failed']:
                continue # Ignore any status but "completed" XXX: catch "progress" too!

            if task_id and jdata['name'] == f'progress/{agent_id}/task/{task_id}':
                break

    return task_id

def run(agent_id, action, data={}, **kwargs):
    """Run a new task and wait until it completes. An object with exit_code, error, output is returned.
    """
    results = runp([{
        'agent_id': agent_id,
        'action': action,
        'data': data,
    }], **kwargs)
    return results[0]

def runp(tasks, **kwargs):
    """Run tasks in parallel and return an array of results.
    """
    return asyncio.run(_runp(tasks, **kwargs))

def runp_brief(tasks, **kwargs):
    """Run tasks in parallel and return the number of failed tasks. Errors are sent to sys.stderr
    """
    results = asyncio.run(_runp(tasks, **kwargs))
    errors = 0
    for idx, result in enumerate(results):
        if isinstance(result, Exception):
            print(agent.SD_WARNING+f"Task {tasks[idx]['action']}@{tasks[idx]['agent_id']} run failed: {result}", file=sys.stderr)
            errors += 1
        elif result['exit_code'] != 0:
            print(agent.SD_WARNING+f"Task {tasks[idx]['action']}@{tasks[idx]['agent_id']} run failed: {repr(result)}", file=sys.stderr)
            errors += 1

    return errors

async def _runp(tasks, **kwargs):
    runners = []
    for task in tasks:
        if not 'data' in task:
            task['data']={}
        otask = asyncio.create_task(_run(task['agent_id'], task['action'], data=task['data'], **kwargs), name=task['action'] + '@' + task['agent_id'])
        runners.append(otask)

    return await asyncio.gather(*runners, return_exceptions=(len(tasks) > 1))

async def _run(agent_id, action, data={}, parent=None, progress_range=None, tls_verify=False, endpoint='http://cluster-leader/cluster-admin', auth_token=None):

    if parent is None:
        parent = os.getenv("AGENT_TASK_ID", "")

    # There are two main tasks that run concurrently. We must ensure the web socket captures any status update
    # of the Task even before it is created. This ensures we can catch the "completed" message and retrieve the
    # Task status from Redis after it becomes available.
    async def run_apiserver():
        nonlocal parent
        theaders = _get_theaders_cache(auth_token)
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

    async def run_redis():
        nonlocal parent
        redis_username = os.getenv('REDIS_USER', os.getenv('AGENT_ID', 'default'))
        redis_password = os.environ['REDIS_PASSWORD'] # Fatal if missing!
        async with aioredis.from_url(
            endpoint,
            username=redis_username,
            password=redis_password,
            decode_responses=True
        ) as rdb:
            task_id = str(uuid.uuid4())
            task_obj = {
                "id": task_id,
                "action": action,
                "data": data,
                "parent": parent
            }

            pubsub = rdb.pubsub()
            await pubsub.subscribe(f'progress/{agent_id}/task/{task_id}')
            lpush_retval = await rdb.lpush(f'{agent_id}/tasks', json.dumps(task_obj))
            if not (lpush_retval >= 1):
                return None

            async for message in pubsub.listen():
                try:
                    jdata = json.loads(message['data'])
                except:
                    continue # Discard unreadable message

                if jdata['status'] in ['completed', 'aborted', 'validation-failed']:
                    break
                else:
                    continue # Ignore any status but "completed" XXX: catch "progress" too!

            await pubsub.close()

            output = await rdb.get(f'{agent_id}/task/{task_id}/output')
            try:
                output = json.loads(output) # Attempt to decode output as a JSON string
            except:
                pass

            exit_code = await rdb.get(f'{agent_id}/task/{task_id}/exit_code')
            if exit_code.isnumeric():
                exit_code = int(exit_code)

            error = await rdb.get(f'{agent_id}/task/{task_id}/error')

        return {
            'id': task_id,
            'output': output,
            'error': error,
            'exit_code': exit_code,
        } 


    if endpoint.startswith("redis://"):
        return await run_redis()
    else:
        return await run_apiserver()
