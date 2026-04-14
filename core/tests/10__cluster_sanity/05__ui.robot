*** Settings ***
Library           Browser

*** Variables ***
${ADMIN_USER}    admin
${ADMIN_PASSWORD}    Nethesis,1234
${MODULE_ID}    ${EMPTY}

*** Keywords ***

Login to cluster-admin
    New Page    https://${NODE_ADDR}/cluster-admin/
    Fill Text    text="Username"    ${ADMIN_USER}
    Click    button >> text="Continue"
    Fill Text    text="Password"    ${ADMIN_PASSWORD}
    Click    button >> text="Log in"
    Wait For Elements State    css=#main-content    visible    timeout=10s

*** Test Cases ***

Take screenshots
    [Tags]    ui
    New Browser    chromium    headless=True
    New Context    ignoreHTTPSErrors=True
    Login to cluster-admin
    Wait For Elements State    h2 >> text="Cluster status"    visible    timeout=10s
    Sleep    5s
    Take Screenshot    filename=${OUTPUT DIR}/browser/screenshot/cluster_status.png
    Go To    https://${NODE_ADDR}/cluster-admin/#/nodes
    Wait For Elements State    h2 >> text="Nodes"    visible    timeout=10s
    Sleep    5s
    Take Screenshot    filename=${OUTPUT DIR}/browser/screenshot/nodes.png
    Go To    https://${NODE_ADDR}/cluster-admin/#/software-center
    Wait For Elements State    h2 >> text="Software center"    visible    timeout=10s
    Sleep    5s
    Take Screenshot    filename=${OUTPUT DIR}/browser/screenshot/software_center.png
    Close Browser
