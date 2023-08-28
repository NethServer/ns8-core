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
import random
import json
import os
import urllib
import sys
from ssl import SSLCertVerificationError
from .exceptions import *

async def run_apiclient_nowait(taskrq, **kwargs):
    kwargs['ssl_ctx'] = None if kwargs.pop('tls_verify', True) is True else False
    theaders = _get_theaders_cache(kwargs.pop('auth_token', None))
    async with aiohttp.ClientSession(raise_for_status=True) as client:
        return await _retry_request(_apost_task, taskrq, client=client, theaders=theaders, **kwargs)

async def run_apiclient(taskrq, **kwargs):
    kwargs['ssl_ctx'] = None if kwargs.pop('tls_verify', True) is True else False
    theaders = _get_theaders_cache(kwargs.pop('auth_token', None))
    async with aiohttp.ClientSession(raise_for_status=True) as client:
        taskctx = {'status_path': None, 'rq': taskrq, 'lock': asyncio.Lock()}

        async with taskctx['lock']:
            tcontroller = asyncio.create_task(_retry_request(_acontrol_task, taskctx, client=client, theaders=theaders, **kwargs), name='tcontroller')
            tpoll = asyncio.create_task(_retry_request(_apoll_status, taskctx, client=client, theaders=theaders, **kwargs), name = 'tpoll')
            try:
                taskctx['status_path'] = await _retry_request(_apost_task, taskrq, client=client, theaders=theaders, **kwargs)
            except:
                tpoll.cancel()
                tcontroller.cancel()
                raise

        def generate_cancel_handler():
            return _cancel_task(
                taskrq['agent_id'],
                taskctx['status_path'].rsplit('/', 1)[1],
                taskrq['parent'],
                kwargs.get('cancel_task_timeout', 0),
                kwargs['endpoint'],
                theaders,
                kwargs['ssl_ctx'],
            )

        # Add a callback to post a cancel-task action:
        kwargs['cancel_handler'].add_callback(generate_cancel_handler)

        done, pending = await asyncio.wait({tcontroller, tpoll}, return_when=asyncio.FIRST_COMPLETED)

        # Remove the cancel callback
        kwargs['cancel_handler'].remove_callback(generate_cancel_handler)

        if tpoll in done:
            tcontroller.cancel()
            return tpoll.result()

        if tcontroller in done:
            tpoll.cancel()
            return tcontroller.result()

async def _alogin(**kwargs):
    """Read the Redis user credentials from the environment and retrieve a new authorization token.
    """
    client = kwargs['client']
    endpoint = kwargs['endpoint']
    redis_username = os.environ['REDIS_USER'] # Fatal if missing!
    redis_password = os.environ['REDIS_PASSWORD'] # Fatal if missing!
    async with client.post(
        f'{endpoint}/api/login',
        json={
            'username':redis_username,
            'password':redis_password
        },
        ssl=kwargs['ssl_ctx'],
    ) as resp:
        jresp = await resp.json()

    if 'AGENT_STATE_DIR' in os.environ:
        cachefile = os.environ['AGENT_STATE_DIR'] + '/apitoken.cache'
        oldmask = os.umask(0o77)
        with open(cachefile, 'w') as fo:
            fo.write(jresp['token'])
        os.umask(oldmask)

    return {"Authorization": "Bearer " + jresp['token']}

async def _aread_status(path, **kwargs):
    """Read the status of the given task_id from agent_id
    """
    client = kwargs['client']
    theaders = kwargs['theaders']
    endpoint = kwargs['endpoint']

    async with client.get(
        f'{endpoint}/api/{path}/status',
        ssl=kwargs['ssl_ctx'],
        headers=theaders,
    ) as resp:
        jresp = await resp.json()

    return {
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

async def _retry_request(request_procedure, *args, **kwargs):
    """Handle request retries and authentication issues"""
    retry_delay = kwargs.pop('retry_initdelay', 1.0)
    retry_delay += random.uniform(0.0, retry_delay)
    attempts = kwargs.pop('retry_attempts', 20)
    incfactor = kwargs.pop('retry_incfactor', 2.0)
    kwargs.setdefault('retry_sendlogin', True)
    http_temporary_errors = [500, 502, 503, 504]
    http_temporary_errors += kwargs.pop('retry_statuscodes', [])

    request_name = request_procedure.__name__
    last_retry_error = ''
    for attempt in range(1, attempts+1):
        try:
            #
            # request_procedure invocation
            #
            retval = await request_procedure(*args, **kwargs)
            if attempt > 1:
                print(f'{request_name} request recovered successfully at attempt {attempt}', file=sys.stderr)
            return retval

        except TaskRetryWsDisconnect as extask:
            retry_error = str(extask)
        except aiohttp.WSServerHandshakeError as exws:
            if exws.status == 200 and kwargs['retry_sendlogin'] is True:
                break # for loop
            else:
                raise exws
        except SSLCertVerificationError as exssl:
            raise exssl
        except aiohttp.ClientConnectorError as exhttp:
            retry_error = str(exhttp)
        except aiohttp.ClientResponseError as exhttp:
            if exhttp.status in http_temporary_errors:
                retry_error = str(exhttp)
            elif exhttp.status == 401 and kwargs['retry_sendlogin'] is True:
                break # for loop. Send a login request before resuming "request_procedure"
            else:
                if attempt > 1:
                    frameback = sys.exc_info()[2].tb_frame.f_back
                    print(f'{request_name} request from {frameback} discarded due to {exhttp} exception at attempt {attempt}', file=sys.stderr)
                raise exhttp
        except Exception as ex:
            if attempt > 1:
                frameback = sys.exc_info()[2].tb_frame.f_back
                print(f'{request_name} request from {frameback} discarded due to {ex} at attempt {attempt}', file=sys.stderr)
            raise

        await asyncio.sleep(retry_delay)
        retry_delay *= incfactor
        if retry_error != last_retry_error:
            print(f'{request_name} request attempt failed ({retry_error}). Retrying...', file=sys.stderr)
            last_retry_error = retry_error
    else:
        raise TaskRetryAbort(f'{request_name} request failed at attempt {attempt}')

    #
    # If the request failed due to a bad auth, retry immediately starting
    # with a login request. Recursion is blocked by the
    # retry_sendlogin=False flag.
    #
    kwargs.pop('retry_sendlogin')
    theaders = await _retry_request(_alogin, retry_sendlogin=False, **kwargs)
    kwargs['theaders'].update(theaders)
    return await _retry_request(request_procedure, *args, retry_sendlogin=False, **kwargs)

async def _cancel_task(agent_id, task_id, parent, timeout, endpoint, theaders, sslctx):
    """Post a cancel-task action
    """
    async with aiohttp.ClientSession(raise_for_status=True) as client:
        async with client.post(
            f'{endpoint}/api/{agent_id}/tasks',
            json={
                'action': 'cancel-task',
                'data': {"task": task_id, "timeout": timeout},
                'parent': parent,
                'extra': {
                    'title': f"{agent_id}/cancel-task",
                    'description': f"Terminate task {task_id}",
                    'isNotificationHidden': True,
                },
            },
            ssl=sslctx,
            headers=theaders,
        ) as resp:

            return await resp.json()

async def _apost_task(taskrq, **kwargs):
    """Create a new task for the given agent_id and return its task_id.
    """
    client = kwargs['client']
    theaders = kwargs['theaders']
    endpoint = kwargs['endpoint']
    agent_id = taskrq['agent_id']

    async with client.post(
        f'{endpoint}/api/{agent_id}/tasks',
        json={
            'action': taskrq['action'],
            'data': taskrq['data'],
            'parent': taskrq['parent'],
            'extra': taskrq.get('extra', {}),
        },
        ssl=kwargs['ssl_ctx'],
        headers=theaders,
    ) as resp:

        jresp = await resp.json()

    return agent_id + '/task/' + jresp['data']['id']

def _get_token(theaders):
    if not 'Authorization' in theaders:
        return ''
    return theaders['Authorization'][len('Bearer '):]

async def _acontrol_task(taskctx, **kwargs):
    client = kwargs['client']
    theaders = kwargs['theaders']
    endpoint = kwargs['endpoint']
    progress_callback = kwargs['progress_callback']
    # Connect the web socket before the task is submitted to catch the
    # relevant progress messages
    async with client.ws_connect(
        f'{endpoint}/ws',
        ssl=kwargs['ssl_ctx'],
        heartbeat=4.0,
    ) as ws:
        await ws.send_json({"action":"authorize","payload":{"jwt":_get_token(theaders)}})
        async for msg in ws:
            if msg.type != aiohttp.WSMsgType.TEXT:
                continue

            try:
                jdata=json.loads(msg.data)
            except:
                continue # Discard unknown message

            async with taskctx['lock']:
                status_path = taskctx['status_path']

            if jdata.get('type') == "task" and jdata.get('name') == 'progress/' + status_path:
                if progress_callback:
                    progress_callback(jdata['payload']['progress'])
                if jdata['payload']['status'] in ['completed', 'aborted', 'validation-failed']:
                    break # for
        else:
            raise TaskRetryWsDisconnect(f"WS reached EOF while waiting for {taskctx['status_path']}")

    return await _retry_request(_aread_status, taskctx['status_path'], retry_statuscodes=[400,404], **kwargs)

async def _apoll_status(taskctx, **kwargs):
    """Read in a loop the status of the given task_id from agent_id
    """
    while True:
        async with taskctx['lock']:
            status_path = taskctx['status_path']

        try:
            return await _aread_status(status_path, **kwargs)
        except aiohttp.ClientResponseError as exhttp:
            if not exhttp.status in [400,404]:
                raise

        await asyncio.sleep(5.0)
