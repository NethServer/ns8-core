*** Settings ***
Library     SSHLibrary
Resource    api.resource

*** Test Cases ***
Get the cluster backup URL
    ${response} =    Run task    cluster/download-cluster-backup    {"password":"0000"}
    Should End With    ${response}[path]    dump.json.gz.gpg
    ${stdout}  ${stderr}  ${rc} =  Execute Command    curl -f --head http://127.0.0.1:9311/backup/${response}[path]    return_stderr=${TRUE}    return_rc=${TRUE}
    Should Be Equal As Integers    ${rc}    0
