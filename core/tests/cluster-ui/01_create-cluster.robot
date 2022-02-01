*** Settings ***
Library           Browser
Suite Setup       Open the browser
Test Setup        Open the context
# *** Variables ***
# //// use variable
# ${NODE_IP}      172.25.5.229 ////

*** Keywords ***
Open the browser
    # New Browser    chromium
    New Browser    chromium    headless=true

Open the context
    New Context    viewport={'width': 1280, 'height': 1024}    ignoreHTTPSErrors=${TRUE}

*** Test Cases ***
Initialize Cluster
    New Page    https://${NODE_ADDR}/cluster-admin/
    Fill Text    input[name="username"]    admin
    Click    button.login-button
    Fill Text    input[name="password"]    Nethesis,1234
    # Fill Text    input[name="password"]    Nethesis,12345
    Click    button.login-button
    Click    text=Create cluster
    Fill Text    text=Current admin password    Nethesis,1234
    Fill Text    text="New admin password"    Nethesis,12345
    Fill Text    text="Re-enter new admin password"    Nethesis,12345
    Click    text=Change password
    # Sleep    1s ////
    # wait until VPN endpoint address is prefilled
    # Wait Until Network Is Idle ////
    # Fill Text    text=VPN endpoint address    ${NODE_ADDR} ////
    Click    button >> text=Create cluster
    # Wait Until Network Is Idle ////
    ${old_browser_timeout} =    Set Browser Timeout    120 seconds
    Log    old browser timeout: ${old_browser_timeout}
    ${old_retry_assertions} =    Set Retry Assertions For    120 seconds
    Log    old retry assertions: ${old_retry_assertions}
    Get Title    ==    Cluster status
    Set Browser Timeout    ${old_browser_timeout}
    Set Retry Assertions For    ${old_retry_assertions}
    # Sleep    10s ////
    # Take Screenshot ////
