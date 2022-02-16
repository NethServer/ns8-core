#
# Copyright (C) 2022 Nethesis S.r.l.
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

import signal

# Global registry of cancel task handlers. Each task
# submissions registers a cancel task handler instance.
# If a TERM/INT signal is received, all handlers run
handler_functions = set()

def _push_cancel_task_handler(loop, handler_function):
    global handler_functions
    if not handler_functions:
        def run_cancel_task_handlers():
            loop.remove_signal_handler(signal.SIGTERM)
            for hf in handler_functions:
                hf()
            handler_functions.clear()
        loop.add_signal_handler(signal.SIGTERM, run_cancel_task_handlers)
    handler_functions.add(handler_function)

def _pop_cancel_task_handler(loop, handler_function):
    global handler_functions
    try:
        handler_functions.remove(handler_function)
    except KeyError:
        pass
    if not handler_functions:
        loop.remove_signal_handler(signal.SIGTERM)
