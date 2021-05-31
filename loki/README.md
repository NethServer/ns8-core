# Loki

Start and configure a instance of Loki, a log aggregation system.
The module use the [Loki official docker image](https://github.com/grafana/loki/releases)

## Install

Instantiate the module:
```
add-module loki 1
```

The output of the command will return the instance name.
Output example:
```
{"rootfull": false, "mid": "loki1"}
```

Wait for `add-module` to complete by looking inside `journalctl`.

## Usage

After the installation, an instance of Loki will be listening on port 3100 in the selected node. The instance can be queried with logcli. eg:
```
root@leader:~# export LOKI_ADDR=http://localhost:3100
root@leader:~# logcli labels -q
__name__
job
nodename
root@leader:~# logcli labels nodename -q
bullseye
leader
root@leader:~# logcli  query -q --no-labels -t '{nodename="leader"} | json | line_format "{{.MESSAGE}}"'
2021-05-28T15:49:27Z Created slice cgroup user-libpod_pod_d32e9a0f2237ca996c90f6790ca90447b3e8cbb30d16cc91701ab8257bb704d6.slice.
2021-05-28T15:49:27Z 2021-05-28 15:49:27.621079036 +0000 UTC m=+0.106351212 container create 6e51520a9fa4f63ac0f3dbbf89ef2ef075041dd90c5df952a35206a14691c654 (image=k8s.gcr.io/pause:3.2, name=d32e9a0f2237-infra)
2021-05-28T15:49:27Z 2021-05-28 15:49:27.62160088 +0000 UTC m=+0.106873062 pod create d32e9a0f2237ca996c90f6790ca90447b3e8cbb30d16cc91701ab8257bb704d6 (image=, name=loki)
2021-05-28T15:49:27Z d32e9a0f2237ca996c90f6790ca90447b3e8cbb30d16cc91701ab8257bb704d6
2021-05-28T15:49:27Z loki.service: Found left-over process 19008 (podman pause) in control group while starting unit. Ignoring.
2021-05-28T15:49:27Z This usually indicates unclean termination of a previous run, or service implementation deficiencies.
```

## Uninstall

To uninstall the instance:
```
remove-module loki1
```
