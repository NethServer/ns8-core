*** Settings ***
Library           SSHLibrary
Suite Teardown    Remove the test command

*** Variables ***
${MODULE}         traefik1

*** Test Cases ***
Module metrics are published
    Write export-metrics    echo robot_test_value 42
    Refresh module metrics
    ${output} =    Get node metrics
    Should Contain    ${output}    robot_test_value{module="${MODULE}"} 42
    Should Contain    ${output}    node_textfile_scrape_error 0

Failed run keeps the last output
    Write export-metrics    exit 1
    Refresh module metrics
    ${output} =    Get node metrics
    Should Contain    ${output}    robot_test_value{module="${MODULE}"} 42

Removed command removes the metrics
    Remove the test command
    ${output} =    Get node metrics
    Should Not Contain    ${output}    robot_test_value
    Should Not Contain    ${output}    module-${MODULE}.prom

*** Keywords ***
Write export-metrics
    [Arguments]    ${body}
    ${rc} =    Execute Command    runagent -m ${MODULE} sh -c 'printf "#!/bin/sh\\n%s\\n" "${body}" > ../bin/export-metrics && chmod +x ../bin/export-metrics'
    ...    return_stdout=False    return_rc=True
    Should Be Equal As Integers    ${rc}    0

Refresh module metrics
    ${rc} =    Execute Command    runagent -m node refresh-module-metrics
    ...    return_stdout=False    return_rc=True
    Should Be Equal As Integers    ${rc}    0

Get node metrics
    ${output}    ${rc} =    Execute Command    curl -sSf http://127.0.0.1:9100/metrics
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0
    RETURN    ${output}

Remove the test command
    Execute Command    runagent -m ${MODULE} rm -f ../bin/export-metrics
    Execute Command    runagent -m node refresh-module-metrics
