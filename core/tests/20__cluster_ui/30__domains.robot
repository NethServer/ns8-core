*** Settings ***
Library           Browser
Resource          ../keywords.resource
Suite Setup       Open the browser
Suite Teardown    Close the browser
Test Setup        Open the context
Test Teardown     Close the context

*** Test Cases ***
Create Internal Domain
    UI login    admin    Nethesis,12345
    Click    text=Domains and users
    Click    button >> text=Create domain
    Click    .ns-tile >> text=Internal
    Click    text=Next
    Click    .cv-modal >> text=Samba
    Click    button >> text=Install provider
    ${old_browser_timeout} =    Set Browser Timeout    60 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    60 seconds
    Fill Text    .cv-text-input >> text="Domain"    test.domain
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
    Fill Text    text="Samba admin password"    Nethesis,1234
    Fill Text    text="Re-enter Samba admin password"    Nethesis,1234
    # Select first IP address
    Click    .bx--list-box__menu-icon
    Click    .bx--list-box__menu-item:first-child
    Click    button >> text=Configure domain
    ${old_browser_timeout} =    Set Browser Timeout    90 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    90 seconds
    Get Element    .cv-tile >> text=test.domain
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
