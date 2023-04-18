#
# Copyright (C) 2023 Nethesis S.r.l.
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

import asyncio

class AsyncSignalHandler:
    """Registry of cancel task handlers. Each task
    submissions registers a cancel task handler instance.
    If a signal is received, all handlers run.

    Usage:

        with AsyncSignalHandler(asyncio.get_running_loop(), signal.SIGTERM) as handler:

            # for some function f()...

            handler.add_callback(f)

            # ...

            handler.remove_callback(f)


    """
    def __init__(self, loop, signal):
        self.signal = signal
        self.handler_generators = set()
        self.loop = loop

    def __enter__(self):
        def run_cancel_task_handlers():
            handlers = [hgen() for hgen in self.handler_generators]
            self.cleanup()
            for hcoro in handlers:
                asyncio.create_task(hcoro) # schedule the handler to run soon

        self.loop.add_signal_handler(self.signal, run_cancel_task_handlers)
        return self

    def __exit__(self, *args):
        self.cleanup()

    def cleanup(self):
        self.handler_generators.clear()
        self.loop.remove_signal_handler(self.signal)

    def add_callback(self, handler_function):
        self.handler_generators.add(handler_function)

    def remove_callback(self, handler_function):
        try:
            self.handler_generators.remove(handler_function)
        except KeyError:
            pass
