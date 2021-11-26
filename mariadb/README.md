# ns8-kickstart

This repository is a starting point to create a NethServer 8 module.

See [NethServer 8 documentation](https://github.com/NethServer/ns8-scratchpad#readme).

Search and replace all "MODULE_" occurrences in this repository with the proper values of your module.

A README.md template for your NethServer 8 module is the following:

# mariadb

Start and configure mariadb instance.
The module uses <MODULE_IMAGE>.

## Install

Instantiate the module:
```
add-module mariadb 1
```

The output of the command will return the instance name.
Output example:
```
{"module_id": "mariadb1", "image_name": "mariadb", "image_url": "<IMAGE_URL>:latest"}
```

## Configure

Let's assume that the mariadb istance is named `mariadb1`.

Then launch `configure-module`, by setting the following parameters:
- `<MODULE_PARAM1_NAME>`: <MODULE_PARAM1_DESCRIPTION>
- `<MODULE_PARAM2_NAME>`: <MODULE_PARAM2_DESCRIPTION>
- ...

Example:
```
[root@fedora mariadb]# api-cli run configure-module --agent module/mariadb7 --data - <<EOF
{ 
"host": "phpmyadmin.domain.com", 
"http2https": true,
"lets_encrypt": false
}
EOF
```

The above command will:
- start and configure the mariadb instance
- (describe configuration process)
- ...

## Uninstall

To uninstall the instance:
```
remove-module mariadb1 --no-preserve
```
# retrieve configuration
```
api-cli run get-configuration --agent module/mariadb12 --data null
Warning: using user "cluster" credentials from the environment
{
"host": "nanatacha.domain.com", 
"http2https": true, 
"lets_encrypt": false
}
```