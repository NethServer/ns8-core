*** Settings ***
Library           Browser
Resource          ../keywords.resource
Suite Setup       Open the browser
Suite Teardown    Close the browser
Test Setup        Open the context
Test Teardown     Close the context

*** Test Cases ***
Create Internal OpenLDAP Domain
    UI login    admin    Nethesis,12345
    Click    text=Domains and users
    Click    button >> text=Create domain
    Click    .ns-tile >> text=Internal
    Click    text=Next
    Click    .cv-modal >> text=OpenLDAP
    Click    button >> text=Install provider
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    ${random} =    Evaluate    random.randint(0, 9999)
    Fill Text    .cv-text-input >> text="Domain"    nethserver.test${random}
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
    Fill Text    text="OpenLDAP admin password"    Nethesis,1234
    Fill Text    text="Re-enter OpenLDAP admin password"    Nethesis,1234
    Click    button >> text=Configure domain
    ${old_browser_timeout} =    Set Browser Timeout    180 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    180 seconds
    Get Element    .cv-tile >> text=nethserver.test${random}
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}

Create Internal Samba Domain
    UI login    admin    Nethesis,12345
    Click    text=Domains and users
    Click    button >> text=Create domain
    Click    .ns-tile >> text=Internal
    Click    text=Next
    Click    .cv-modal >> text=Samba
    Click    button >> text=Install provider
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    ${random} =    Evaluate    random.randint(0, 9999)
    Fill Text    .cv-text-input >> text="Domain"    test.domain${random}
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
    Fill Text    text="Samba admin password"    Nethesis,1234
    Fill Text    text="Re-enter Samba admin password"    Nethesis,1234
    # Select first IP address
    Click    .bx--list-box__menu-icon
    Click    .bx--list-box__menu-item:first-child
    Click    button >> text=Configure domain
    ${old_browser_timeout} =    Set Browser Timeout    180 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    180 seconds
    Get Element    .cv-tile >> text=test.domain${random}
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
