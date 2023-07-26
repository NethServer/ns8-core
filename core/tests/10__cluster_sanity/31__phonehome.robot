*** Settings ***
Library     SSHLibrary
Resource    api.resource


*** Test Cases ***
Check if phonehome.timer is enabled
    SSHLibrary.File Should Exist    /etc/systemd/system/timers.target.wants/phonehome.timer

Check correct output from get-facts
    ${response} =    Run task    cluster/get-facts    null
    Should Not Be Empty    ${response['$schema']}
    Should Not Be Empty    ${response['uuid']}
    Should Be Equal    ${response['installation']}    nethserver
    Should Not Be Empty    ${response['facts']}
