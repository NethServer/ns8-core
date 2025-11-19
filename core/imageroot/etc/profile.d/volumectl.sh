
#
# Copyright (C) 2025 Nethesis S.r.l.
# SPDX-License-Identifier: GPL-3.0-or-later
#

if [ -z "$BASH" ]; then
    return
fi

# Bash completion for volumectl
_volumectl_complete()
{
    local cur prev cmd
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    cmd="${COMP_WORDS[1]}"

    # Extract app names (sections) from /etc/nethserver/volumes.conf
    local apps=()
    if [[ -r /etc/nethserver/volumes.conf ]]; then
        apps=( $(grep '^\[' /etc/nethserver/volumes.conf | sed 's/^\[\(.*\)\]$/\1/' | sort -u) )
    fi

    # Extract base mount paths from volumectl list-base-paths
    local paths=()
    local volumes=()
    if command -v volumectl >/dev/null 2>&1; then
        paths=( $(volumectl list-base-paths 2>/dev/null | awk '{print $1}') )
        volumes=( $(volumectl list-volumes 2>/dev/null) )
    fi

    if [[ $COMP_CWORD -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "list-base-paths add-volume remove-volume" -- "$cur") )
        return
    fi

    case "$cmd" in
        list-base-paths)
            # No options
            ;;
        add-volume)
            case "$prev" in
                --target)
                    COMPREPLY=( $(compgen -W "${paths[*]}" -- "$cur") )
                    ;;
                --for)
                    COMPREPLY=( $(compgen -W "${apps[*]}" -- "$cur") )
                    ;;
                *)
                    COMPREPLY=( $(compgen -W "--target --for --next-only" -- "$cur") )
                    ;;
            esac
            ;;
        remove-volume)
            case "$prev" in
                --for)
                    COMPREPLY=( $(compgen -W "${apps[*]}" -- "$cur") )
                    ;;
                *)
                    # Suggest volume names from config, plus ALL
                    COMPREPLY=( $(compgen -W "${volumes[*]} --for" -- "$cur") )
                    ;;
            esac
            ;;
    esac
}

complete -F _volumectl_complete volumectl
