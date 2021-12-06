# mariadb && PhpMyAdmin

Start and configure mariadb instance.

The module uses : 

docker.io/library/phpmyadmin  5.1.1
docker.io/library/mariadb     10.6.5

## Install

Instantiate the module:

```
add-module ghcr.io/nethserver/mariadb:latest 1
```

The output of the command will return the instance name.
Output example:
```
{'module_id': 'mariadb1', 'image_name': 'mariadb', 'image_url': 'ghcr.io/nethserver/mariadb:latest'}
```

## Configure

Let's assume that the mariadb istance is named `mariadb1`.

Then launch `configure-module`, by setting the following parameters:

Example:
```
[root@fedora mariadb]# api-cli run configure-module --agent module/mariadb7 --data - <<EOF
{ 
"path": "/phpmyadmin, 
"http2https": true
}
EOF
```

The above command will:
- validate the /webpath (if not already used by another application.)
- start and configure the mariadb instance
- retrieve two random ports, one for the phpmyadmin container the other for the external port to connect the mysql client


## Uninstall

To uninstall the instance:

```
remove-module mariadb1 --no-preserve
```

## retrieve configuration

You can retrieve the configuration like the user interface 

```
api-cli run get-configuration --agent module/mariadb12 --data null
Warning: using user "cluster" credentials from the environment
{
"path": "/path", 
"http2https": true
}
```