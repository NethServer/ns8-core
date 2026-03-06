*** Settings ***
Library     SSHLibrary
Resource    api.resource


*** Test Cases ***
Create user with allowed_networks
    Run task    cluster/add-user
    ...    {"user":"testnet","password_hash":"73cb3858a687a8494ca3323053016282f3dad39d42cf62ca4e79dda2aac7d9ac","set":{"display_name":"Test Net","allowed_networks":["192.168.1.0/24","10.0.0.1"]},"grant":[]}

get-user-info returns allowed_networks
    ${jinfo} =    Run task    cluster/get-user-info    {"user":"testnet"}
    Should Contain    ${jinfo}[allowed_networks]    192.168.1.0/24
    Should Contain    ${jinfo}[allowed_networks]    10.0.0.1

alter-user updates allowed_networks
    Run task    cluster/alter-user
    ...    {"user":"testnet","set":{"allowed_networks":["172.16.0.0/12"]}}
    ${jinfo} =    Run task    cluster/get-user-info    {"user":"testnet"}
    Should Contain    ${jinfo}[allowed_networks]    172.16.0.0/12
    Should Not Contain Any    ${jinfo}[allowed_networks]    192.168.1.0/24    10.0.0.1

alter-user clears allowed_networks
    Run task    cluster/alter-user    {"user":"testnet","set":{"allowed_networks":[]}}
    ${jinfo} =    Run task    cluster/get-user-info    {"user":"testnet"}
    Should Be Empty    ${jinfo}[allowed_networks]

Remove test user
    Run task    cluster/remove-user    {"user":"testnet"}
