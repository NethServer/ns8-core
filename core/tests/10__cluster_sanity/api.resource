*** Settings ***
Library    SSHLibrary

*** Keywords ***
Run task
    [Arguments]    ${action}    ${input}    ${decode_json}=${TRUE}    ${rc_expected}=0
    ${stdout}    ${stderr}    ${rc} =     Execute Command    api-cli run ${action} --data '${input}'    return_stdout=True    return_stderr=True    return_rc=True
    Should Be Equal As Integers    ${rc_expected}    ${rc}    Run task ${action} failed!${\n}${stderr}
    IF    ${decode_json} and len($stdout) > 0
        ${response} =    Evaluate    json.loads('''${stdout}''')    modules=json
    ELSE
        ${response} =    Set Variable    ${stdout}
    END
    [Return]    ${response}
