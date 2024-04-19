
#
# Copyright (C) 2024 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

# PATH override for login shell -- fixes #6888 Override PATH when this
# script is evaluated under the agent environment (runagent):
if [ -n "${AGENT_STATE_DIR}" ] && [ -f "${AGENT_STATE_DIR}/agent.env" ]; then
   . ${AGENT_STATE_DIR}/agent.env
   export PATH
fi
