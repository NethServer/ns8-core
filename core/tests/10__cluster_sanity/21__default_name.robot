*** Settings ***
Library     SSHLibrary
Resource    api.resource


*** Test Cases ***
Check if ui_name is generated
    ${response} =    Run task    cluster/get-name    null
    Should Not Be Empty    ${response['name']}
