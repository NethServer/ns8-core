#
# Import the agent environment to run NethServer commands manually.
# This is useful during development. We should remove it later.
#
set -a
source /etc/nethserver/agent.env
set +a

# Declare pathmunge function if not defined
# The function is already available inside Fedora profile
declare -F pathmunge &>/dev/null || pathmunge ()
{
    case ":${PATH}:" in
        *:"$1":*)

        ;;
        *)
            if [ "$2" = "after" ]; then
                PATH=$PATH:$1;
            else
                PATH=$1:$PATH;
            fi
        ;;
    esac
}

if [[ -d ~/.config/bin ]]; then
    # Push rootless image binary path
    pathmunge ~/.config/bin
fi

if [[ -r ~/.config/state/agent.env ]]; then
    set -a
    source ~/.config/state/agent.env
    set +a
fi

if [[ -r /var/lib/nethserver/cluster/state/agent.env ]]; then
    set -a
    source /var/lib/nethserver/cluster/state/agent.env
    set +a
fi

alias logcli="/var/lib/nethserver/node/logcli.sh"
