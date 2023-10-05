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

# we want to be compatible python 3.9 and 3.11, aioredis is inside redis 4.5.5
try :
    from redis import asyncio as aioredis
except:
    import aioredis

import asyncio
import json
import os
import sys
import uuid
import random
from .exceptions import *

async def _task_submission_check_client_idle(rdb, taskrq, max_idle_time):
    for client_item in await rdb.client_list('normal'):
        if client_item['name'] == taskrq['agent_id']:
            if int(client_item['idle']) <= max_idle_time:
                break
    else:
        raise TaskSubmissionCheckFailed(f"Client \"{taskrq['agent_id']}\" was not found")

async def run_redisclient_nowait(taskrq, **kwargs):
    redis_username = os.environ['REDIS_USER'] # Fatal if missing!
    redis_password = os.environ['REDIS_PASSWORD'] # Fatal if missing!
    async with aioredis.from_url(
        kwargs['endpoint'],
        username=redis_username,
        password=redis_password,
        decode_responses=True
    ) as rdb:
        if kwargs['check_idle_time'] > 0:
            await _task_submission_check_client_idle(rdb, taskrq, kwargs['check_idle_time'])

        task_id = str(uuid.uuid4())
        task_obj = {
            "id": task_id,
            "action": taskrq['action'],
            "data": taskrq['data'],
            "parent": taskrq['parent'],
            'extra': taskrq.get('extra', {}),
        }
        await _retry_request(rdb.lpush, f'{taskrq["agent_id"]}/tasks', json.dumps(task_obj))

    return f"task/{taskrq['agent_id']}/{task_id}"

async def run_redisclient(taskrq, **kwargs):
    redis_username = os.environ['REDIS_USER'] # Fatal if missing!
    redis_password = os.environ['REDIS_PASSWORD'] # Fatal if missing!
    async with aioredis.from_url(
        kwargs['endpoint'],
        username=redis_username,
        password=redis_password,
        decode_responses=True,
    ) as rdb:
        if kwargs['check_idle_time'] > 0:
            await _task_submission_check_client_idle(rdb, taskrq, kwargs['check_idle_time'])

        task_id = str(uuid.uuid4())
        taskctx = {
            'id': task_id,
            'status_path': f"task/{taskrq['agent_id']}/{task_id}",
            'progress_channel': f"progress/{taskrq['agent_id']}/task/{task_id}",
            'rq': taskrq,
            'lock': asyncio.Lock(),
            'pushed': False,
        }

        def generate_cancel_handler():
            return _cancel_task(
                kwargs['endpoint'],
                redis_username,
                redis_password,
                taskrq['agent_id'],
                task_id,
                taskrq['parent'],
                kwargs.get('cancel_task_timeout', 0),
            )

        # Add a callback to post a cancel-task action:
        kwargs['cancel_handler'].add_callback(generate_cancel_handler)

        # Context is initially locked. Unlock occurs when task has been
        # submitted by _acontrol_task:
        await taskctx['lock'].acquire()

        tcontrol = asyncio.create_task(_retry_request(_acontrol_task, rdb, taskctx, **kwargs), name='tcontrol')
        tpoll = asyncio.create_task(_retry_request(_apoll_status, rdb, taskctx), name='tpoll')

        done, pending = await asyncio.wait({tcontrol, tpoll}, return_when=asyncio.FIRST_COMPLETED)

        # Remove the cancel callback
        kwargs['cancel_handler'].remove_callback(generate_cancel_handler)

        if tpoll in done:
            tcontrol.cancel()
            return tpoll.result()

        if tcontrol in done:
            tpoll.cancel()
            return tcontrol.result()

async def _cancel_task(
        endpoint,
        username,
        password,
        agent_id,
        task_id,
        parent_task,
        timeout,
    ):
    """Push a cancel-task action
    """
    async with aioredis.from_url(
        endpoint,
        username=username,
        password=password,
        decode_responses=True,
    ) as rdb:
        print(f'Sending cancel-task request to task/{agent_id}/{task_id}', file=sys.stderr)
        return await rdb.lpush(f"{agent_id}/tasks", json.dumps({
            'id': str(uuid.uuid4()),
            'action': 'cancel-task',
            'data': {'task': task_id, 'timeout': timeout},
            'parent': parent_task,
            'extra': {
                'title': f"{agent_id}/cancel-task",
                'description': f"Terminate task {task_id}",
                'isNotificationHidden': True,
            },
        }))

async def _apoll_status(rdb, taskctx):
    while True:
        async with taskctx['lock']:
            try:
                return await _aread_status(rdb, taskctx['status_path'])
            except TaskStatusNotFound:
                pass

        await asyncio.sleep(5.0)

async def _aread_status(rdb, status_path):

    output, error, exit_code = await rdb.mget(
        status_path + '/output',
        status_path + '/error',
        status_path + '/exit_code',
    )

    if output is None or exit_code is None or error is None:
        raise TaskStatusNotFound('Task status %s not found' % status_path)

    try:
        output = json.loads(output) # Attempt to decode output as a JSON string
    except:
        pass

    return {
        'output': output,
        'error': error,
        'exit_code': int(exit_code),
    }

async def _acontrol_task(rdb, taskctx, **kwargs):
    progress_callback = kwargs['progress_callback']
    async with rdb.pubsub() as pubsub:
        await _retry_request(pubsub.subscribe, taskctx['progress_channel'])

        if taskctx['lock'].locked() and taskctx['pushed'] is False:
            taskrq = taskctx['rq']
            task_obj = {
                "id": taskctx['id'],
                "action": taskrq["action"],
                "data": taskrq["data"],
                "parent": taskrq["parent"],
                "extra": taskrq.get("extra", {}),
            }
            await _retry_request(rdb.lpush, taskrq['agent_id'] + '/tasks', json.dumps(task_obj))
            taskctx['pushed'] = True
            taskctx['lock'].release() # _apoll_status is waiting on the lock

        async for message in pubsub.listen():
            try:
                jdata = json.loads(message['data'])
            except:
                continue # Discard unreadable message

            if progress_callback:
                progress_callback(jdata['progress'])

            if jdata['status'] in ['completed', 'aborted', 'validation-failed']:
                break
        else:
            raise TaskRetryPubSubDisconnect()

        return await _retry_request(_aread_status, rdb, taskctx['status_path'])

async def _retry_request(request_procedure, *args, **kwargs):
    """Handle request retries and authentication issues"""
    retry_delay = kwargs.pop('retry_initdelay', 1.0)
    retry_delay += random.uniform(0.0, retry_delay)
    attempts = kwargs.pop('retry_attempts', 20)
    incfactor = kwargs.pop('retry_incfactor', 1.2)

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
        except aioredis.ConnectionError as exredis:
            retry_error = str(exredis)
        except TaskRetryPubSubDisconnect as extask:
            retry_error = str(exredis)
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
