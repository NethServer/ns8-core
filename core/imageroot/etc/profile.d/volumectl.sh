
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
        COMPREPLY=( $(compgen -W "list-base-paths add-volume remove-volume move-volume" -- "$cur") )
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
        move-volume)
            case "$prev" in
                --target)
                    local default_path
                    default_path=$(findmnt -n -o TARGET --target "${HOME_BASEDIR:-/home}" 2>/dev/null)
                    COMPREPLY=( $(compgen -W "${paths[*]} ${default_path}" -- "$cur") )
                    ;;
                --for)
                    # Rootless application instances of this node
                    COMPREPLY=( $(compgen -W "$(ls /var/lib/systemd/linger 2>/dev/null)" -- "$cur") )
                    ;;
                *)
                    # Volume names of the application given with --for
                    local i mid=""
                    for (( i=2; i < COMP_CWORD; i++ )); do
                        if [[ ${COMP_WORDS[i]} == --for ]]; then
                            mid=${COMP_WORDS[i+1]}
                        fi
                    done
                    local modvolumes=""
                    if [[ -n $mid ]]; then
                        modvolumes=$(runagent -m "$mid" podman volume ls -q 2>/dev/null)
                    fi
                    COMPREPLY=( $(compgen -W "${modvolumes} --for --target --sync --force" -- "$cur") )
                    ;;
            esac
            ;;
    esac
}

complete -F _volumectl_complete volumectl
