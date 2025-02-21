
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

if [ -z "$BASH" ]; then
    return
fi

_runagent_wastyped()
{
    local xword checkword="$1"
    for xword in "${COMP_WORDS[@]}"; do
        if [[ "${xword}" == "${checkword}" ]]; then
            return 0
        fi
    done
    return 1
}

_runagent_completions()
{
    local cword="$2"
    local agents=()

    if _runagent_wastyped "-l" || _runagent_wastyped "--list-modules"; then
        :
    elif _runagent_wastyped "-m" || _runagent_wastyped "--module-id" ; then
        agents+=($(runagent --list-modules | sort | grep ^"${cword}"))
        COMPREPLY+=($(compgen -W "${agents[*]}" -- "${cword}"))
    else
        COMPREPLY+=($(compgen -W "--list-modules --module-id --current-dir --help" -- "${cword}"))
    fi
}

complete -F _runagent_completions runagent
