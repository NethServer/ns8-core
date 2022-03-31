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
    ${old_browser_timeout} =    Set Browser Timeout    240 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    240 seconds
    Get Text    .bx--toast-notification--success >> text=installed on
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}

Install Another Instance
    UI login    admin    Nethesis,12345
    Click    text="Software center"
    Click    .app-list .bx--row > div:first-child button >> text="Instances"
    Click    button >> text="Install new instance"
    Click    .bx--btn--primary >> text="Install"
    ${old_browser_timeout} =    Set Browser Timeout    240 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    240 seconds
    Get Text    .bx--toast-notification--success >> text=installed on
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}

Clone App
    UI login    admin    Nethesis,12345
    Click    text="Software center"
    Click    .app-list .bx--row > div:first-child button >> text="Instances"
    # click overflow menu
    Click    .bx--row > div:first-child .cv-overflow-menu[data-test-id="first"] button
    Click    .bx--overflow-menu-options [data-test-id="first-clone"] >> text="Clone"
    Click    button >> text="Clone app"
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    Get Text    .bx--toast-notification--success >> text=cloned to
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}

Uninstall App
    UI login    admin    Nethesis,12345
    Click    text="Software center"
    Click    .app-list .bx--row > div:first-child button >> text="Instances"
    # click overflow menu
    Click    .bx--row > div:first-child .cv-overflow-menu[data-test-id="first"] button
    Click    .bx--overflow-menu-options [data-test-id="first-uninstall"] >> text="Uninstall"
    # get module ID from modal title
    ${modal_title}=    Get Text    .bx--modal-header__heading >> text=Uninstall
    ${regexp_match}=    Evaluate    re.search("Uninstall (.+)", "${modal_title}"), re
    # enter module ID in danger modal input
    Fill Text    .cv-modal .bx--text-input    ${regexp_match[0].group(1)}
    Click    button >> text="I understand, delete"
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    Get Text    .bx--toast-notification--success >> text="Completed"
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
