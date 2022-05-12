#
# Copyright (C) 2021 Nethesis S.r.l.
# http://www.nethesis.it - nethserver@nethesis.it
#
# This script is part of NethServer.
#
# NethServer is free software: you can redisclienttribute it and/or modify
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

import asyncio
import os
import sys
import signal

from .apiclient import run_apiclient, run_apiclient_nowait
from .redisclient import run_redisclient, run_redisclient_nowait
from .exceptions import *
from .handlers import AsyncSignalHandler

def run(agent_id, action, data={}, **kwargs):
    """Run a new task and wait until it completes. An object with exit_code, error, output is returned.
    """
    taskrq = {
        'agent_id': agent_id,
        'action': action,
        'data': data,
    }
    results = runp([taskrq], **kwargs)
    return results[0]

def run_nowait(agent_id, action, data={}, **kwargs):
    """Run the task but do no wait for completion. Instead return immediately a task identifier.
    """
    return run(agent_id, action, data, nowait=True, **kwargs)

def runp(tasks, **kwargs):
    """Run tasks in parallel and return an array of results.
    """
    return asyncio.run(_runp(tasks, **kwargs))

def runp_nowait(tasks, **kwargs):
    """Run tasks in parallel but do not wait for completion. Instead
    return immediately an array of task identifiers.
    """
    return asyncio.run(_runp(tasks, nowait=True, **kwargs))

def runp_brief(tasks, **kwargs):
    """Run tasks in parallel and return the number of failed tasks. Errors are sent to sys.stderr
    """
    results = asyncio.run(_runp(tasks, **kwargs))
    errors = 0
    for idx, result in enumerate(results):
        if isinstance(result, Exception):
            print(f"Task {tasks[idx]['agent_id']}/{tasks[idx]['action']} run error: {result}", file=sys.stderr)
            errors += 1
        elif result['exit_code'] != 0:
            print(f"Task {tasks[idx]['agent_id']}/{tasks[idx]['action']} run failed: {repr(result)}", file=sys.stderr)
            errors += 1

    return errors

async def _runp(tasks, **kwargs):

    if 'progress_callback' in kwargs and not kwargs['progress_callback'] is None:
        # Equally distribute the progress weight of each task
        parent_cbk = kwargs['progress_callback']
        del kwargs['progress_callback']
        runp_progress = [0] * len(tasks)
        last_value = -1
        def create_task_cbk(idx):
            return lambda p: task_progress_callback(p, idx)
        def task_progress_callback(p, idx):
            nonlocal last_value, runp_progress, parent_cbk
            runp_progress[idx] = p
            curr_value = int(sum(runp_progress) / len(runp_progress))
            if curr_value > last_value:
                last_value = curr_value
                parent_cbk(curr_value)
    else:
        parent_cbk = None

    nowait = kwargs.pop('nowait', False)
    kwargs.setdefault('check_idle_time', 0 if nowait else 8) # Check the client connection is alive

    with AsyncSignalHandler(asyncio.get_running_loop(), signal.SIGTERM) as cancel_handler:
        runners = []
        for idx, taskrq in enumerate(tasks):
            if not 'data' in taskrq:
                taskrq['data']={}

            if not 'parent' in taskrq:
                taskrq['parent'] = os.getenv("AGENT_TASK_ID", "")

            if 'extra' in kwargs:
                taskrq['extra'] = kwargs['extra']

            if parent_cbk:
                task_cbk = create_task_cbk(idx)
            else:
                task_cbk = None

            if nowait:
                tcoro = _run_with_protocol_nowait(taskrq, **kwargs)
            else:
                tcoro = _run_with_protocol(taskrq, progress_callback=task_cbk, cancel_handler=cancel_handler, **kwargs)

            runners.append(asyncio.create_task(tcoro, name=taskrq['agent_id'] + '/' + taskrq['action']))

        return await asyncio.gather(*runners, return_exceptions=(len(tasks) > 1))


async def _run_with_protocol(taskrq, **pconn):
    pconn.setdefault('progress_callback', None)
    pconn.setdefault('endpoint', 'http://cluster-leader:9311')
    if pconn['endpoint'].startswith("redis://"):
        return await run_redisclient(taskrq, **pconn)
    else:
        return await run_apiclient(taskrq, **pconn)

async def _run_with_protocol_nowait(taskrq, **pconn):
    pconn.setdefault('endpoint', 'http://cluster-leader:9311')
    if pconn['endpoint'].startswith("redis://"):
        return await run_redisclient_nowait(taskrq, **pconn)
    else:
        return await run_apiclient_nowait(taskrq, **pconn)
