---
layout: default
title: Migration
nav_order: 8
---

# Migration from NS7

**THIS IS A WORK-IN-PROGRESS PROCEDURE, DO NOT USE IT IN A PRODUCTION ENVIRONMENT**

Install the NS8 cluster, make sure to set a strong password for `admin` user, then make sure to request
a valid certificate for the server FQDN:
```
api-cli run set-certificate --agent module/traefik1 --data "{\"fqdn\": \"$(hostname -f)\"}"
```

To verify if the certificate has been correctly generated, use:
```
api-cli run get-certificate --agent module/traefik1 --data "{\"fqdn\": \"$(hostname -f)\"}"
```
The output should be something like:
```
{"fqdn": "ns8n1.nethserver.org", "obtained": true}
```


1. Install the migration package on NS7:
   ```
   yum --enablerepo=nethserver-testing install nethserver-ns8-migration
   ```

2. On NS7 machine, join the NS8 cluster:
   ```
   ns8-join ns8n1.nethserver.org admin Nethesis,1234
   ```

3. Apply the configuration and setup the VPN to the cluster:
   ```
   signal-event nethserver-ns8-migration-save
   ```

## Nextcloud migration

If NS7 has a running AD, go to NS7 and execute:
```
/usr/share/nethesis/nethserver-ns8-migration/apps/ldapproxy/export
```

If everything is fine, the output should look like this:
```
{"code": 200, "data": {"error": "", "exit_code": 0, "file": "cluster/task/d2b90404-79c6-4b72-8886-9d5b3e804954", "output": true}, "message": "success"}
{"code": 200, "data": {"error": "", "exit_code": 0, "file": "cluster/task/fd7131a5-6d7e-491b-8b6b-d541f4e68131", "output": ""}, "message": "success"}
```

Then on NS7 execute:
```
/usr/share/nethesis/nethserver-ns8-migration/apps/nextcloud/export
```

The above command will install Nextcloud to NS8 machine and start the first synchronization.
The command can be executed multiple times.
When you are ready to migrate, put the Nextcloud instance in maintenance mode, then execute:
```
/usr/share/nethesis/nethserver-ns8-migration/apps/nextcloud/migrate
```

Please note that above command will fail if Nextcloud is not configured to be accessed using a FQDN.
