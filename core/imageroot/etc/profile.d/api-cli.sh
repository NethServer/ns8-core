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

if [ -z "$BASH" ]; then
    return
fi

_apicli_wastyped()
{
    local xword checkword="$1"
    for xword in "${COMP_WORDS[@]}"; do
        if [[ "${xword}" == "${checkword}" ]]; then
            return 0
        fi
    done
    return 1
}

_apicli_completions()
{
    local cword="$2"
    local actions=()

    if _apicli_wastyped "run"; then
        actions+=($(api-cli list-actions "${cword}"))
        COMPREPLY+=($(compgen -W "${actions[*]}" -- "${cword}"))
    elif _apicli_wastyped "login"; then
        COMPREPLY+=($(compgen -W "--username --password --output --help" -- "${cword}"))
    elif _apicli_wastyped "logout"; then
        :
    else
        COMPREPLY+=($(compgen -W "login logout run --help" -- "${cword}"))
    fi
}

complete -F _apicli_completions api-cli
