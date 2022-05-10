*** Settings ***
Library           Browser
Resource          ../keywords.resource
Suite Setup       Open the browser
Suite Teardown    Close the browser
Test Setup        Open the context
Test Teardown     Close the context

*** Test Cases ***
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
    Fill Text    .cv-text-input >> text="Domain"    nethserver.sambatest
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
    Get Element    .cv-tile >> text=nethserver.sambatest
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}

Create Samba user
    UI login    admin    Nethesis,12345
    Create domain user    nethserver.sambatest    aaatestuser    Test User

Edit Samba user
    UI login    admin    Nethesis,12345
    Edit domain user    nethserver.sambatest    aaatestuser    Test User Edited

Delete Samba user
    UI login    admin    Nethesis,12345
    Delete domain user    nethserver.sambatest    aaatestuser

Create Samba group
    UI login    admin    Nethesis,12345
    Create domain group    nethserver.sambatest    aaatestgroup    Test group description

Edit Samba group
    UI login    admin    Nethesis,12345
    Edit domain group    nethserver.sambatest    aaatestgroup    Test group description edited

Delete Samba group
    UI login    admin    Nethesis,12345
    Delete domain group    nethserver.sambatest    aaatestgroup

Delete Internal Samba Domain
    UI login    admin    Nethesis,12345
    Click    text=Domains and users
    # click overflow menu
    Click    .cv-overflow-menu[data-test-id="nethserver.sambatest-menu"] button
    Click    .bx--overflow-menu-options [data-test-id="nethserver.sambatest-delete"] >> text="Delete"
    # enter domain name in danger modal input
    Fill Text    .cv-modal.bx--modal--danger .bx--text-input    nethserver.sambatest
    Click    button >> text="I understand, delete"
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    Get Text    .bx--toast-notification--success >> text="Completed"
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
