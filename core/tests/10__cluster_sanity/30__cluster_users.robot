*** Settings ***
Library     SSHLibrary
Resource    api.resource


*** Test Cases ***
Admin user is listed
    ${lusers} =    Run task    cluster/list-users    ""    decode_json=${FALSE}
    Should Contain    ${lusers}    "admin"

Set admin user display_name
    Run task    cluster/alter-user    {"user":"admin","set":{"display_name":"Admin"}}

Check admin user display_name
    ${jinfo} =    Run task    cluster/get-user-info    {"user":"admin"}
    Should Be Equal As Strings    ${jinfo}[display_name]    Admin
