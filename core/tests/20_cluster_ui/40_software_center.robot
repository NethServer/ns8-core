*** Settings ***
Library           Browser
Resource          ../keywords.resource
Suite Setup       Open the browser
Suite Teardown    Close the browser
Test Setup        Open the context
Test Teardown     Close the context

*** Test Cases ***
Install App
    UI login    admin    Nethesis,12345
    Click    text="Software center"
    Click    .app-list .bx--row > div:first-child button >> text="Install"
    Click    .bx--btn--primary >> text="Install"
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    Get Text    .bx--toast-notification--success >> text=installed on
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}

Install Another Instance
    UI login    admin    Nethesis,12345
    Click    text="Software center"
    Click    .app-list .bx--row > div:first-child button >> text="Instances"
    Click    button >> text="Install new instance"
    Click    .bx--btn--primary >> text="Install"
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    Get Text    .bx--toast-notification--success >> text=installed on
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}

Uninstall App
    UI login    admin    Nethesis,12345
    Click    text="Software center"
    Click    .app-list .bx--row > div:first-child button >> text="Instances"
    # click overflow menu
    Click    .bx--row > div:first-child .cv-overflow-menu button
    # temporarily disable strict mode (query can return more than one element)
    ${old_mode} =    Set Strict Mode    False
    Click    .bx--overflow-menu-options >> text="Uninstall"
    Set Strict Mode    ${old_mode}
    Click    button >> text="Uninstall instance"
    Get Text    .bx--toast-notification--success >> text="Completed"
