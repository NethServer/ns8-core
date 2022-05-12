*** Settings ***
Resource    api.resource

*** Variables ***
${mid1}
${domain}

*** Keywords ***
Domain is ready
    Wait Until Keyword Succeeds    10    1 seconds
    ...    Run task    cluster/list-domain-users    {"domain":"${domain}"}

Create users and groups
    Run task    module/${mid1}/add-user    {"user":"u1","display_name":"First User","password":"Nethesis,1234"}
    Run task    module/${mid1}/add-user    {"user":"u2","display_name":"Second User","locked":true}
    Run task    module/${mid1}/add-user    {"user":"u3","display_name":"Third User","groups":["domain admins"],"password":"Nethesis,1234"}
    Run task    module/${mid1}/add-group   {"group":"g1","description":"Group One","users":["u1"]}

List users
    ${response} =     Run task    cluster/list-domain-users    {"domain":"${domain}"}    decode_json=${FALSE}
    Should Contain    ${response}    "u1"
    Should Contain    ${response}    "First User"
    Should Match Regexp    ${response}    "locked": *(false|true)
    Should Contain    ${response}    "u2"
    Should Contain    ${response}    "Second User"
    Should Contain    ${response}    "u3"
    Should Contain    ${response}    "Third User"
    Should Contain Any   ${response}    "Administrator"    "admin"    "administrator"

List groups
    ${response} =     Run task    cluster/list-domain-groups    {"domain":"${domain}"}    decode_json=${FALSE}
    Should Contain    ${response}    "g1"
    Should Contain    ${response}    "Group One"
    Should Match Regexp    ${response}    "users": *\\["u1"\\]

Check First User is unlocked
    ${response} =     Run task    cluster/get-domain-user    {"domain":"${domain}","user":"u1"}
    Should Be Equal    ${response['user']['user']}    u1
    Should Be Equal    ${response['user']['display_name']}    First User
    Should Be Equal    ${response['user']['locked']}    ${FALSE}

Check First User is member of g1
    ${response} =     Run task    cluster/get-domain-user    {"domain":"${domain}","user":"u1"}
    Should Be Equal    ${response['user']['groups'][0]['group']}    g1
    Should Be Equal    ${response['user']['groups'][0]['description']}    Group One

Check Second User is locked
    ${response} =     Run task    cluster/get-domain-user    {"domain":"${domain}","user":"u2"}
    Should Be Equal    ${response['user']['user']}    u2
    Should Be Equal    ${response['user']['display_name']}    Second User
    Should Be Equal    ${response['user']['locked']}    ${TRUE}

Check Group One members
    ${response} =     Run task    cluster/get-domain-group    {"domain":"${domain}","group":"g1"}
    Should Be Equal    ${response['group']['group']}    g1
    Should Be Equal    ${response['group']['description']}    Group One
    Should Be Equal    ${response['group']['users'][0]['user']}    u1
    Should Be Equal    ${response['group']['users'][0]['display_name']}    First User

Get a non-existing user
    ${response} =    Run task    cluster/get-domain-user    {"domain":"${domain}","user":"UUU"}    rc_expected=2

Get a non-existing group
    ${response} =    Run task    cluster/get-domain-group    {"domain":"${domain}","group":"GGG"}    rc_expected=2
