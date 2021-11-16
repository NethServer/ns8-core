# ns8-kickstart

This repository is a starting point to create a NethServer 8 module.

See [NethServer 8 documentation](https://github.com/NethServer/ns8-scratchpad#readme).

Search and replace all "MODULE_" occurrences in this repository with the proper values of your module.

A README.md template for your NethServer 8 module is the following:

# <MODULE_NAME>

Start and configure <MODULE_NAME> instance.
The module uses <MODULE_IMAGE>.

## Install

Instantiate the module:
```
add-module <MODULE_NAME> 1
```

The output of the command will return the instance name.
Output example:
```
{"module_id": "<MODULE_NAME>1", "image_name": "<MODULE_NAME>", "image_url": "<IMAGE_URL>:latest"}
```

## Configure

Let's assume that the <MODULE_NAME> istance is named `<MODULE_NAME>1`.

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
- start and configure the <MODULE_NAME> instance
- (describe configuration process)
- ...

## Uninstall

To uninstall the instance:
```
remove-module <MODULE_NAME>1 --no-preserve
```
