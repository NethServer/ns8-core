---
layout: default
title: Database
nav_order: 90
parent: Modules
---

# Database access

Everyone can access the Redis database with read-only privileges. The
Redis user `default` (password is not important) can be used for this
purpose.

To get write access a module must provide the Redis credentials stored in
its `agent.env` file. The complete path is `~/.config/state/agent.env`.
Write access is restricted to Redis keys and channels with prefix
`module/{module_id}/*`. The same credentials allow broad read access, also
on `private/agents/*` namespace.

The above rules are already implemented by the Python `agent` module.

Access Redis in read-only:
```python
import agent

rdb = agent.redis_connect(privileged=False)
somehash = rdb.hgetall('cluster/somehash')
```

Access Redis in read-write mode:
```python
import agent

rdb = agent.redis_connect(privileged=True)
rdb.hset('module/myapp1/myhash', mapping={'myvar': 'myvalue'})
```

If Redis connection is required at service boot time, prefer connecting to
the local replica. This avoids issues if the leader node is unreachable.
```python
import agent

rdb = agent.redis_connect(use_replica=True)
cluster_network = rdb.get('cluster/network')
```
