*** Settings ***
Library           Browser
Resource          ../keywords.resource
Suite Setup       Open the browser
Suite Teardown    Close the browser
Test Setup        Open the context
Test Teardown     Close the context

*** Test Cases ***
Dump Logs
    UI login    admin    Nethesis,12345
    Click    text="System logs"
    Click    button >> text="Search"
    # look for a Gin log
    Get Text    .logs-output >> text=GIN

Follow Logs
    UI login    admin    Nethesis,12345
    Click    text="System logs"
    Click    .cv-content-switcher >> text="Follow"
    Click    button.bx--btn--primary >> text="Follow"
    # look for a Gin log
    Get Text    .logs-output >> text=GIN
