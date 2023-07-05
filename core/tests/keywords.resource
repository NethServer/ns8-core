*** Settings ***
Documentation     This is a resource file, that can contain variables and keywords.
...               Keywords defined here can be used where this Keywords.resource in loaded.
Library           SSHLibrary
Library           DateTime
Library           Browser

*** Variables ***
${SSH_KEYFILE}    %{HOME}/.ssh/id_ecdsa
${HEADLESS}       true

*** Keywords ***
Connect to the node
    [Arguments]    ${NODE_ADDR}    ${SSH_KEYFILE}
    Open Connection    ${NODE_ADDR}
    Login With Public Key    root    ${SSH_KEYFILE}    keep_alive_interval=30 seconds

System boot is complete
    ${output} =    Execute Command    systemctl is-system-running    --wait
    Should Be True    '${output}' == 'running' or '${output}' == 'degraded'

Install the core
    [Arguments]    ${COREMODULES}
    Save the journal begin timestamp
    ${COREMODULES} =    Replace String    ${COREMODULES}    ,    ${SPACE}
    Put File    install.sh    .    mode=0644
    ${rc} =    Execute Command    bash install.sh ${COREMODULES}
    ...    return_stdout=False
    ...    return_stderr=False
    ...    return_rc=True
    Should Be Equal As Integers    ${rc}    0

Save the journal begin timestamp
    ${tsnow} =    Get Current Date    result_format=epoch
    Set Global Variable    ${JOURNAL_SINCE}    ${tsnow}

Basic post install sanity checks
    SSHLibrary.File Should Exist    /etc/nethserver/wg0.key
    SSHLibrary.File Should Exist    /etc/nethserver/wg0.pub
    SSHLibrary.File Should Exist    /etc/nethserver/api-server.env

Uninstall core and modules
    Put File    imageroot/var/lib/nethserver/node/uninstall.sh    .    mode=0644
    Execute Command    bash uninstall.sh
    Collect the suite journal

Collect the suite journal
    Execute Command    journalctl -S @${JOURNAL_SINCE} >journal-dump.log
    Get File    journal-dump.log    ${OUTPUT DIR}/journal-${SUITE NAME}.log

Open the browser
    ${old_browser_timeout} =    Set Browser Timeout    20 seconds
    ${old_retry_assertions} =    Set Retry Assertions For    20 seconds
    # Log To Console    old_browser_timeout ${old_browser_timeout}
    # Log To Console    old_retry_assertions ${old_retry_assertions}
    New Browser    chromium    headless=${HEADLESS}

Open the context
    New Context    viewport={'width': 1280, 'height': 1024}    ignoreHTTPSErrors=${TRUE}

Close the context
    Close Context    CURRENT

Close the browser
    Close Browser    CURRENT

UI login
    [Arguments]    ${username}    ${password}
    Open login page
    Fill Text    text="Username"    ${username}
    Click    button >> text="Continue"
    Fill Text    text="Password"    ${password}
    Click    button >> text="Log in"

Open login page
    # support for localhost:PORT as NODE_ADDR
    IF    $NODE_ADDR.startswith("localhost")
        ${url}=    Set Variable    http://${NODE_ADDR}
    ELSE
        ${url}=    Set Variable    https://${NODE_ADDR}/cluster-admin
    END
    New Page    ${url}

Create domain user
    [Arguments]    ${domain}    ${username}    ${display_name}
    Click    text=Domains and users
    Click    button[data-test-id="${domain}-users-and-groups"]
    Click    .toolbar button >> text=Create user
    Fill Text    .cv-text-input >> text="Username"    ${username}
    Sleep    200ms
    Fill Text    .cv-text-input >> text="Display name"    ${display_name}
    Sleep    200ms
    Fill Text    .cv-text-input >> text="Password"    Nethesis,1234
    Sleep    200ms
    Fill Text    .cv-text-input >> text="Re-enter password"    Nethesis,1234
    Click    .cv-modal button >> text=Create user
    Get Text    .bx--toast-notification--success >> text="Completed"
    Get Element    .ns-data-table td >> text="${username}"

Edit domain user
    [Arguments]    ${domain}    ${username}    ${new_display_name}
    Click    text=Domains and users
    Click    button[data-test-id="${domain}-users-and-groups"]
    # click overflow menu
    Click    .cv-overflow-menu[data-test-id="${username}-menu"] button
    Click    .bx--overflow-menu-options [data-test-id="${username}-edit"] >> text="Edit"
    Fill Text    .cv-text-input >> text="Display name"    ${new_display_name}
    Click    .cv-modal button >> text=Edit user
    Get Text    .bx--toast-notification--success >> text="Completed"
    Get Element    .ns-data-table td >> text="${new_display_name}"

Delete domain user
    [Arguments]    ${domain}    ${username}
    Click    text=Domains and users
    Click    button[data-test-id="${domain}-users-and-groups"]
    # click overflow menu
    Click    .cv-overflow-menu[data-test-id="${username}-menu"] button
    Click    .bx--overflow-menu-options [data-test-id="${username}-delete"] >> text="Delete"
    # enter domain name in danger modal input
    Fill Text    [data-test-id="delete-user-modal"] .bx--text-input    ${username}
    Click    [data-test-id="delete-user-modal"] button >> text="I understand, delete"
    Get Text    .bx--toast-notification--success >> text="Completed"

Create domain group
    [Arguments]    ${domain}    ${group}    ${description}
    Click    text=Domains and users
    Click    button[data-test-id="${domain}-users-and-groups"]
    Click    .toolbar button >> text=Create group
    Fill Text    .cv-text-input >> text="Name"    ${group}
    Sleep    200ms
    Fill Text    .cv-text-input >> text="Description"    ${description}
    Click    .cv-modal button >> text=Create group
    Get Text    .bx--toast-notification--success >> text="Completed"
    Get Element    .ns-data-table td >> text="${group}"

Edit domain group
    [Arguments]    ${domain}    ${group}    ${new_description}
    Click    text=Domains and users
    Click    button[data-test-id="${domain}-users-and-groups"]
    # click overflow menu
    Click    .cv-overflow-menu[data-test-id="${group}-menu"] button
    Click    .bx--overflow-menu-options [data-test-id="${group}-edit"] >> text="Edit"
    Fill Text    .cv-text-input >> text="Description"    ${new_description}
    Click    .cv-modal button >> text=Edit group
    Get Text    .bx--toast-notification--success >> text="Completed"
    Get Element    .ns-data-table td >> text="${new_description}"

Delete domain group
    [Arguments]    ${domain}    ${group}
    Click    text=Domains and users
    Click    button[data-test-id="${domain}-users-and-groups"]
    # click overflow menu
    Click    .cv-overflow-menu[data-test-id="${group}-menu"] button
    Click    .bx--overflow-menu-options [data-test-id="${group}-delete"] >> text="Delete"
    # enter domain name in danger modal input
    Fill Text    [data-test-id="delete-group-modal"] .bx--text-input    ${group}
    Click    [data-test-id="delete-group-modal"] button >> text="I understand, delete"
    Get Text    .bx--toast-notification--success >> text="Completed"
