*** Settings ***
Library           Browser
Resource          ../keywords.resource
Suite Setup       Open the browser
Suite Teardown    Close the browser
Test Setup        Open the context
Test Teardown     Close the context

*** Test Cases ***
Create Cluster
    [Tags]    install
    UI login    admin    Nethesis,1234
    Click    text=Create cluster
    Fill Text    text=Current admin password    Nethesis,1234
    Fill Text    text="New admin password"    Nethesis,12345
    Fill Text    text="Re-enter new admin password"    Nethesis,12345
    Click    button >> text=Change password
    Click    button >> text=Create cluster
    ${old_browser_timeout} =    Set Browser Timeout    240 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    240 seconds
    Get Title    ==    Cluster status
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
