#
# Import the agent environment to run NethServer commands manually.
# This is useful during development. We should remove it later.
#
set -a
source /etc/nethserver/agent.env
set +a

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
